package genorm

import (
	"context"
	"errors"
	"fmt"
)

type UpdateContext[T Table] struct {
	*Context[T]
	assignExprs    []*TableAssignExpr[T]
	whereCondition whereConditionClause[T]
	order          orderClause[T]
	limit          limitClause
}

func Update[T Table](table T, options ...Option) *UpdateContext[T] {
	ctx := newContext(table, options...)

	return &UpdateContext[T]{
		Context: ctx,
	}
}

func (c *UpdateContext[T]) Set(assignExprs ...*TableAssignExpr[T]) (res *UpdateContext[T]) {
	if len(assignExprs) == 0 {
		c.addError(errors.New("no assign expressions"))
		return c
	}

	c.assignExprs = append(c.assignExprs, assignExprs...)

	return c
}

func (c *UpdateContext[T]) Where(
	condition TypedTableExpr[T, WrappedPrimitive[bool]],
) *UpdateContext[T] {
	err := c.whereCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("where condition: %w", err))
	}

	return c
}

func (c *UpdateContext[T]) OrderBy(direction OrderDirection, expr TableExpr[T]) *UpdateContext[T] {
	err := c.order.add(orderItem[T]{
		expr:      expr,
		direction: direction,
	})
	if err != nil {
		c.addError(fmt.Errorf("order by: %w", err))
	}

	return c
}

func (c *UpdateContext[T]) Limit(limit uint64) *UpdateContext[T] {
	err := c.limit.set(limit)
	if err != nil {
		c.addError(fmt.Errorf("limit: %w", err))
	}

	return c
}

func (c *UpdateContext[T]) DoCtx(ctx context.Context, db DB) (rowsAffected int64, err error) {
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

func (c *UpdateContext[T]) Do(db DB) (rowsAffected int64, err error) {
	return c.DoCtx(context.Background(), db)
}

func (c *UpdateContext[T]) buildQuery() (string, []ExprType, error) {
	args := []ExprType{}
	qb := c.config.newQueryBuilder()

	str := "UPDATE "
	err := qb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	tableQuery, tableArgs, errs := c.table.Expr()
	if len(errs) != 0 {
		return "", nil, fmt.Errorf("table expr: %w", errs[0])
	}

	err = qb.WriteExprQuery(tableQuery)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", tableQuery, err)
	}

	args = append(args, tableArgs...)

	if len(c.assignExprs) == 0 {
		return "", nil, errors.New("no assignment")
	}

	str = " SET "
	err = qb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	for i, expr := range c.assignExprs {
		if i > 0 {
			err = qb.WriteString(", ")
			if err != nil {
				return "", nil, fmt.Errorf("write string: %w", err)
			}
		}

		assignmentQuery, assignmentArgs, errs := expr.AssignExpr()
		if len(errs) != 0 {
			return "", nil, errs[0]
		}

		err = qb.WriteExprQuery(assignmentQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write assignment: %w", err)
		}
		args = append(args, assignmentArgs...)
	}

	if c.whereCondition.exists() {
		whereQuery, whereArgs, err := c.whereCondition.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("where condition: %w", err)
		}

		str = " WHERE "
		err = qb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		err = qb.WriteExprQuery(whereQuery)
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
		err = qb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		err = qb.WriteExprQuery(orderQuery)
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
		err = qb.WriteString(str)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		err = qb.WriteExprQuery(limitQuery)
		if err != nil {
			return "", nil, fmt.Errorf("write string(%s): %w", limitQuery, err)
		}

		args = append(args, limitArgs...)
	}

	return qb.String(), args, nil
}
