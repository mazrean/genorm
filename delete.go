package genorm

import (
	"context"
	"fmt"
	"strings"
)

type DeleteContext[T BasicTable] struct {
	*Context[T]
	whereCondition whereConditionClause[T]
	order          orderClause[T]
	limit          limitClause
}

func Delete[T BasicTable](table T) *DeleteContext[T] {
	ctx := newContext(table)

	return &DeleteContext[T]{
		Context: ctx,
	}
}

func (c *DeleteContext[T]) Where(
	condition TypedTableExpr[T, WrappedPrimitive[bool]],
) *DeleteContext[T] {
	err := c.whereCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("where condition: %w", err))
	}

	return c
}

func (c *DeleteContext[T]) OrderBy(direction OrderDirection, expr TableExpr[T]) *DeleteContext[T] {
	err := c.order.add(orderItem[T]{
		expr:      expr,
		direction: direction,
	})
	if err != nil {
		c.addError(fmt.Errorf("order by: %w", err))
	}

	return c
}

func (c *DeleteContext[T]) Limit(limit uint64) *DeleteContext[T] {
	err := c.limit.set(limit)
	if err != nil {
		c.addError(fmt.Errorf("limit: %w", err))
	}

	return c
}

func (c *DeleteContext[T]) DoCtx(ctx context.Context, db DB) (rowsAffected int64, err error) {
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

func (c *DeleteContext[T]) Do(db DB) (rowsAffected int64, err error) {
	return c.DoCtx(context.Background(), db)
}

func (c *DeleteContext[T]) buildQuery() (string, []ExprType, error) {
	args := []ExprType{}

	sb := strings.Builder{}

	str := "DELETE FROM `"
	_, err := sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	str = c.table.TableName()
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	str = "`"
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
