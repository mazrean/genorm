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

func (c *DeleteContext[TableBase]) Where(condition genorm.TypedTableExpr[TableBase, bool]) *DeleteContext[TableBase] {
	err := c.whereCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("where condition: %w", err))
	}

	return c
}

func (c *DeleteContext[TableBase]) OrderBy(direction OrderDirection, expr genorm.TableExpr[TableBase]) *DeleteContext[TableBase] {
	err := c.order.add(orderItem[TableBase]{
		expr:      expr,
		direction: direction,
	})
	if err != nil {
		c.addError(fmt.Errorf("order by: %w", err))
	}

	return c
}

func (c *DeleteContext[TableBase]) Limit(limit uint64) *DeleteContext[TableBase] {
	err := c.limit.set(limit)
	if err != nil {
		c.addError(fmt.Errorf("limit: %w", err))
	}

	return c
}

func (c *DeleteContext[TableBase]) DoContext(ctx context.Context, db *sql.DB) (rowsAffected int64, err error) {
	errs := c.Errors()
	if len(errs) != 0 {
		return 0, errs[0]
	}

	query, args, err := c.buildQuery()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
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

func (c *DeleteContext[TableBase]) Do(db *sql.DB) (rowsAffected int64, err error) {
	return c.DoContext(context.Background(), db)
}

func (c *DeleteContext[BasicTable]) buildQuery() (string, []any, error) {
	args := []any{}

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
