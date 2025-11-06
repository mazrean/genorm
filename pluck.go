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
	groupExpr       groupClause[T]
	havingCondition whereConditionClause[T]
	order           orderClause[T]
	limit           limitClause
	offset          offsetClause
	lockType        lockClause
}

func Pluck[T Table, S ExprType](table T, field TypedTableExpr[T, S], options ...Option) *PluckContext[T, S] {
	return &PluckContext[T, S]{
		Context: newContext(table, options...),
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
	err := c.groupExpr.set(exprs)
	if err != nil {
		c.addError(fmt.Errorf("group by: %w", err))
	}

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
	err := c.lockType.set(lockType)
	if err != nil {
		c.addError(fmt.Errorf("lock: %w", err))
	}

	return c
}

func (c *PluckContext[T, S]) GetAllCtx(ctx context.Context, db DB) ([]S, error) {
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

	query = c.config.formatQuery(query)

	rows, err := db.QueryContext(ctx, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return []S{}, nil
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

func (c *PluckContext[T, S]) GetAll(db DB) ([]S, error) {
	return c.GetAllCtx(context.Background(), db)
}

func (c *PluckContext[T, S]) GetCtx(ctx context.Context, db DB) (S, error) {
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

	query = c.config.formatQuery(query)

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

func (c *PluckContext[T, S]) Get(db DB) (S, error) {
	return c.GetCtx(context.Background(), db)
}

func (c *PluckContext[T, S]) buildQuery() (string, []ExprType, error) {
	sb := strings.Builder{}
	args := []ExprType{}

	str := "SELECT "
	_, err := sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write select(%s): %w", str, err)
	}

	if c.distinct {
		str = "DISTINCT "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write distinct(%s): %w", str, err)
		}
	}

	fieldQuery, fieldArgs, errs := c.field.Expr()
	if len(errs) != 0 {
		return "", nil, fmt.Errorf("field: %w", errs[0])
	}

	_, err = sb.WriteString(fieldQuery)
	if err != nil {
		return "", nil, fmt.Errorf("write field(%s): %w", fieldQuery, err)
	}

	str = " AS res"
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write as(%s): %w", str, err)
	}

	args = append(args, fieldArgs...)

	str = " FROM "
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write from(%s): %w", str, err)
	}

	tableQuery, tableArgs, errs := c.table.Expr()
	if len(errs) != 0 {
		return "", nil, fmt.Errorf("table expr: %w", errs[0])
	}

	_, err = sb.WriteString(tableQuery)
	if err != nil {
		return "", nil, fmt.Errorf("write table(%s): %w", tableQuery, err)
	}

	args = append(args, tableArgs...)

	if c.whereCondition.exists() {
		whereQuery, whereArgs, err := c.whereCondition.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("where condition: %w", err)
		}

		str = " WHERE "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write where(%s): %w", str, err)
		}

		_, err = sb.WriteString(whereQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write where(%s): %w", whereQuery, err)
		}

		args = append(args, whereArgs...)
	}

	if c.groupExpr.exists() {
		groupExpr, groupArgs, err := c.groupExpr.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("group expr: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(groupExpr)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", groupExpr, err)
		}

		args = append(args, groupArgs...)
	}

	if c.havingCondition.exists() {
		havingQuery, havingArgs, err := c.havingCondition.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("having condition: %w", err)
		}

		str = " HAVING "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write having(%s): %w", str, err)
		}

		_, err = sb.WriteString(havingQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write having(%s): %w", havingQuery, err)
		}

		args = append(args, havingArgs...)
	}

	if c.order.exists() {
		orderQuery, orderArgs, err := c.order.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("order: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write order(%s): %w", str, err)
		}

		_, err = sb.WriteString(orderQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write order(%s): %w", orderQuery, err)
		}

		args = append(args, orderArgs...)
	}

	if c.limit.exists() {
		limitQuery, limitArgs, err := c.limit.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("limit: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write limit(%s): %w", str, err)
		}

		_, err = sb.WriteString(limitQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write limit(%s): %w", limitQuery, err)
		}

		args = append(args, limitArgs...)
	}

	if c.offset.exists() {
		offsetQuery, offsetArgs, err := c.offset.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("offset: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write offset(%s): %w", str, err)
		}

		_, err = sb.WriteString(offsetQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write offset(%s): %w", offsetQuery, err)
		}

		args = append(args, offsetArgs...)
	}

	if c.lockType.exists() {
		lockQuery, lockArgs, err := c.lockType.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("lock type: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write lock(%s): %w", str, err)
		}

		_, err = sb.WriteString(lockQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write lock(%s): %w", lockQuery, err)
		}

		args = append(args, lockArgs...)
	}

	return sb.String(), args, nil
}
