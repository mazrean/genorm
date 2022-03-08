package genorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type PluckContext[T Table, S ExprType] struct {
	*Context[T]
	distinct        bool
	field           TypedTableExpr[T, S]
	whereCondition  whereConditionClause[T]
	groupExpr       []TableExpr[T]
	havingCondition whereConditionClause[T]
	order           orderClause[T]
	limit           limitClause
	offset          offsetClause
	lockType        LockType
}

func Pluck[T Table, S ExprType](table T, field TypedTableExpr[T, S]) *PluckContext[T, S] {
	return &PluckContext[T, S]{
		Context: newContext(table),
		field:   field,
	}
}

func (c *PluckContext[T, S]) Distinct() *PluckContext[T, S] {
	if c.distinct {
		c.addError(errors.New("distinct already set"))
		return c
	}

	c.distinct = true

	return c
}

func (c *PluckContext[T, S]) Where(
	condition TypedTableExpr[T, WrappedPrimitive[bool]],
) *PluckContext[T, S] {
	err := c.whereCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("where condition: %w", err))
	}

	return c
}

func (c *PluckContext[T, S]) GroupBy(exprs ...TableExpr[T]) *PluckContext[T, S] {
	if len(exprs) == 0 {
		c.addError(errors.New("no group expr"))
		return c
	}

	c.groupExpr = append(c.groupExpr, exprs...)

	return c
}

func (c *PluckContext[T, S]) Having(
	condition TypedTableExpr[T, WrappedPrimitive[bool]],
) *PluckContext[T, S] {
	err := c.havingCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("having condition: %w", err))
	}

	return c
}

func (c *PluckContext[T, S]) OrderBy(direction OrderDirection, expr TableExpr[T]) *PluckContext[T, S] {
	err := c.order.add(orderItem[T]{
		expr:      expr,
		direction: direction,
	})
	if err != nil {
		c.addError(fmt.Errorf("order by: %w", err))
	}

	return c
}

func (c *PluckContext[T, S]) Limit(limit uint64) *PluckContext[T, S] {
	err := c.limit.set(limit)
	if err != nil {
		c.addError(fmt.Errorf("limit: %w", err))
	}

	return c
}

func (c *PluckContext[T, S]) Offset(offset uint64) *PluckContext[T, S] {
	err := c.offset.set(offset)
	if err != nil {
		c.addError(fmt.Errorf("offset: %w", err))
	}

	return c
}

func (c *PluckContext[T, S]) Lock(lockType LockType) *PluckContext[T, S] {
	if c.lockType != none {
		c.addError(errors.New("lock already set"))
		return c
	}

	if lockType != ForUpdate && lockType != ForShare {
		c.addError(errors.New("invalid lock type"))
		return c
	}

	c.lockType = lockType

	return c
}

func (c *PluckContext[T, S]) FindCtx(ctx context.Context, db DB) ([]S, error) {
	errs := c.Errors()
	if len(errs) != 0 {
		return nil, errs[0]
	}

	query, exprArgs, err := c.buildQuery()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	args := make([]any, 0, len(exprArgs))
	for _, arg := range exprArgs {
		args = append(args, arg)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	exprs := []S{}
	for rows.Next() {
		var expr S

		err := rows.Scan(&expr)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		exprs = append(exprs, expr)
	}

	return exprs, nil
}

func (c *PluckContext[T, S]) Find(db DB) ([]S, error) {
	return c.FindCtx(context.Background(), db)
}

func (c *PluckContext[T, S]) TakeCtx(ctx context.Context, db DB) (S, error) {
	var res S

	err := c.limit.set(1)
	if err != nil {
		return res, fmt.Errorf("set limit 1: %w", err)
	}

	errs := c.Errors()
	if len(errs) != 0 {
		return res, errs[0]
	}

	query, queryArgs, err := c.buildQuery()
	if err != nil {
		return res, fmt.Errorf("build query: %w", err)
	}

	args := make([]any, 0, len(queryArgs))
	for _, arg := range queryArgs {
		args = append(args, arg)
	}

	row := db.QueryRowContext(ctx, query, args...)

	err = row.Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return res, ErrRecordNotFound
	}
	if err != nil {
		return res, fmt.Errorf("query: %w", err)
	}

	return res, nil
}

func (c *PluckContext[T, S]) Take(db DB) (S, error) {
	return c.TakeCtx(context.Background(), db)
}

func (c *PluckContext[T, S]) buildQuery() (string, []ExprType, error) {
	sb := strings.Builder{}
	args := []ExprType{}

	sb.WriteString("SELECT ")

	if c.distinct {
		sb.WriteString("DISTINCT ")
	}

	fieldQuery, fieldArgs, errs := c.field.Expr()
	if len(errs) != 0 {
		return "", nil, fmt.Errorf("field: %w", errs[0])
	}
	sb.WriteString(fieldQuery)
	sb.WriteString(" AS res")
	args = append(args, fieldArgs...)

	sb.WriteString(" FROM ")

	tableQuery, tableArgs, errs := c.table.Expr()
	if len(errs) != 0 {
		return "", nil, fmt.Errorf("table expr: %w", errs[0])
	}

	sb.WriteString(tableQuery)
	args = append(args, tableArgs...)

	if c.whereCondition.exists() {
		whereQuery, whereArgs, err := c.whereCondition.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("where condition: %w", err)
		}

		sb.WriteString(" WHERE ")
		sb.WriteString(whereQuery)
		args = append(args, whereArgs...)
	}

	if len(c.groupExpr) != 0 {
		sb.WriteString(" GROUP BY ")

		groupQueries := make([]string, 0, len(c.groupExpr))
		for _, groupExpr := range c.groupExpr {
			groupQuery, groupArgs, errs := groupExpr.Expr()
			if len(errs) != 0 {
				return "", nil, errs[0]
			}

			groupQueries = append(groupQueries, groupQuery)
			args = append(args, groupArgs...)
		}

		sb.WriteString(strings.Join(groupQueries, ", "))
	}

	if c.havingCondition.exists() {
		havingQuery, havingArgs, err := c.havingCondition.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("having condition: %w", err)
		}

		sb.WriteString(" HAVING ")
		sb.WriteString(havingQuery)
		args = append(args, havingArgs...)
	}

	if c.order.exists() {
		orderQuery, orderArgs, err := c.order.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("order: %w", err)
		}

		sb.WriteString(" ")
		sb.WriteString(orderQuery)
		args = append(args, orderArgs...)
	}

	if c.limit.exists() {
		limitQuery, limitArgs, err := c.limit.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("limit: %w", err)
		}

		sb.WriteString(" ")
		sb.WriteString(limitQuery)
		args = append(args, limitArgs...)
	}

	if c.offset.exists() {
		offsetQuery, offsetArgs, err := c.offset.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("offset: %w", err)
		}

		sb.WriteString(" ")
		sb.WriteString(offsetQuery)
		args = append(args, offsetArgs...)
	}

	if c.lockType != none {
		switch c.lockType {
		case ForUpdate:
			sb.WriteString(" FOR UPDATE")
		case ForShare:
			sb.WriteString(" FOR SHARE")
		}
	}

	return sb.String(), args, nil
}
