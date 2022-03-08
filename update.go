package genorm

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type UpdateContext[T Table] struct {
	*Context[T]
	assignExprs    []*TableAssignExpr[T]
	whereCondition whereConditionClause[T]
	order          orderClause[T]
	limit          limitClause
}

func Update[T Table](table T) *UpdateContext[T] {
	ctx := newContext(table)

	return &UpdateContext[T]{
		Context: ctx,
	}
}

func (c *UpdateContext[Table]) Set(assignExprs ...*TableAssignExpr[Table]) (res *UpdateContext[Table]) {
	if len(assignExprs) == 0 {
		c.addError(errors.New("no assign expressions"))
		return c
	}

	c.assignExprs = append(c.assignExprs, assignExprs...)

	return c
}

func (c *UpdateContext[Table]) Where(
	condition TypedTableExpr[Table, WrappedPrimitive[bool]],
) *UpdateContext[Table] {
	err := c.whereCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("where condition: %w", err))
	}

	return c
}

func (c *UpdateContext[Table]) OrderBy(direction OrderDirection, expr TableExpr[Table]) *UpdateContext[Table] {
	err := c.order.add(orderItem[Table]{
		expr:      expr,
		direction: direction,
	})
	if err != nil {
		c.addError(fmt.Errorf("order by: %w", err))
	}

	return c
}

func (c *UpdateContext[Table]) Limit(limit uint64) *UpdateContext[Table] {
	err := c.limit.set(limit)
	if err != nil {
		c.addError(fmt.Errorf("limit: %w", err))
	}

	return c
}

func (c *UpdateContext[Table]) DoCtx(ctx context.Context, db DB) (rowsAffected int64, err error) {
	errs := c.Errors()
	if len(errs) != 0 {
		return 0, errs[0]
	}

	query, exprArgs, err := c.buildQuery()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	args := make([]any, 0, len(exprArgs))
	for _, arg := range exprArgs {
		args = append(args, arg)
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (c *UpdateContext[Table]) Do(db DB) (rowsAffected int64, err error) {
	return c.DoCtx(context.Background(), db)
}

func (c *UpdateContext[Table]) buildQuery() (string, []ExprType, error) {
	args := []ExprType{}

	sb := strings.Builder{}

	str := "UPDATE "
	_, err := sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	tableQuery, tableArgs, errs := c.table.Expr()
	if len(errs) != 0 {
		return "", nil, fmt.Errorf("table expr: %w", errs[0])
	}

	_, err = sb.WriteString(tableQuery)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", tableQuery, err)
	}

	args = append(args, tableArgs...)

	if len(c.assignExprs) == 0 {
		return "", nil, errors.New("no assignment")
	}

	str = " SET "
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	assignments := make([]string, 0, len(c.assignExprs))
	for _, expr := range c.assignExprs {
		assignmentQuery, assignmentArgs, errs := expr.AssignExpr()
		if len(errs) != 0 {
			return "", nil, errs[0]
		}

		assignments = append(assignments, assignmentQuery)
		args = append(args, assignmentArgs...)
	}

	str = strings.Join(assignments, ", ")
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	if c.whereCondition.exists() {
		whereQuery, whereArgs, err := c.whereCondition.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("where condition: %w", err)
		}

		str = " WHERE "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(whereQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", whereQuery, err)
		}

		args = append(args, whereArgs...)
	}

	if c.order.exists() {
		orderQuery, orderArgs, err := c.order.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("order: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(orderQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", orderQuery, err)
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
			return "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(limitQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", limitQuery, err)
		}

		args = append(args, limitArgs...)
	}

	return sb.String(), args, nil
}
