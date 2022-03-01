package statement

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/mazrean/genorm"
)

type DeleteContext[T BasicTable] struct {
	*Context[T]
	whereCondition whereConditionClause[T]
	order          orderClause[T]
	limit          limitClause
}

func NewDeleteContext[T BasicTable](table T) *DeleteContext[T] {
	ctx := newContext(table)

	return &DeleteContext[T]{
		Context: ctx,
	}
}

func (c *DeleteContext[Table]) Where(
	condition genorm.TypedTableExpr[Table, genorm.WrappedPrimitive[bool]],
) *DeleteContext[Table] {
	err := c.whereCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("where condition: %w", err))
	}

	return c
}

func (c *DeleteContext[Table]) OrderBy(direction OrderDirection, expr genorm.TableExpr[Table]) *DeleteContext[Table] {
	err := c.order.add(orderItem[Table]{
		expr:      expr,
		direction: direction,
	})
	if err != nil {
		c.addError(fmt.Errorf("order by: %w", err))
	}

	return c
}

func (c *DeleteContext[Table]) Limit(limit uint64) *DeleteContext[Table] {
	err := c.limit.set(limit)
	if err != nil {
		c.addError(fmt.Errorf("limit: %w", err))
	}

	return c
}

func (c *DeleteContext[Table]) DoContext(ctx context.Context, db *sql.DB) (rowsAffected int64, err error) {
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

func (c *DeleteContext[Table]) Do(db *sql.DB) (rowsAffected int64, err error) {
	return c.DoContext(context.Background(), db)
}

func (c *DeleteContext[Table]) buildQuery() (string, []genorm.ExprType, error) {
	args := []genorm.ExprType{}

	sb := strings.Builder{}
	sb.WriteString("DELETE FROM ")
	sb.WriteString(c.table.TableName())

	if c.whereCondition.exists() {
		whereQuery, whereArgs, err := c.whereCondition.getExpr()
		if err != nil {
			return "", nil, fmt.Errorf("where condition: %w", err)
		}

		sb.WriteString(" WHERE ")
		sb.WriteString(whereQuery)
		args = append(args, whereArgs...)
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

	return sb.String(), args, nil
}
