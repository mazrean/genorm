package statement

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mazrean/genorm"
)

type UpdateContext[T Table] struct {
	*Context[T]
	assignExprs    []*genorm.TableAssignExpr[T]
	whereCondition whereConditionClause[T]
	order          orderClause[T]
	limit          limitClause
}

func NewUpdateContext[T Table](table T) *UpdateContext[T] {
	ctx := newContext(table)

	return &UpdateContext[T]{
		Context: ctx,
	}
}

func (c *UpdateContext[Table]) Set(assignExprs ...*genorm.TableAssignExpr[Table]) (res *UpdateContext[Table]) {
	if len(assignExprs) == 0 {
		c.addError(errors.New("no assign expressions"))
		return c
	}

	c.assignExprs = append(c.assignExprs, assignExprs...)

	return c
}

func (c *UpdateContext[Table]) Where(
	condition genorm.TypedTableExpr[Table, genorm.WrappedPrimitive[bool]],
) *UpdateContext[Table] {
	err := c.whereCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("where condition: %w", err))
	}

	return c
}

func (c *UpdateContext[Table]) OrderBy(direction OrderDirection, expr genorm.TableExpr[Table]) *UpdateContext[Table] {
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

func (c *UpdateContext[Table]) DoContext(ctx context.Context, db *sql.DB) (rowsAffected int64, err error) {
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

func (c *UpdateContext[Table]) Do(db *sql.DB) (rowsAffected int64, err error) {
	return c.DoContext(context.Background(), db)
}

func (c *UpdateContext[Table]) buildQuery() (string, []genorm.ExprType, error) {
	args := []genorm.ExprType{}

	sb := strings.Builder{}
	sb.WriteString("UPDATE ")

	tableQuery, tableArgs, errs := c.table.Expr()
	if len(errs) != 0 {
		return "", nil, fmt.Errorf("table expr: %w", errs[0])
	}

	sb.WriteString(tableQuery)
	args = append(args, tableArgs...)

	if len(c.assignExprs) == 0 {
		return "", nil, errors.New("no assignment")
	}

	sb.WriteString(" SET ")
	assignments := make([]string, 0, len(c.assignExprs))
	for _, expr := range c.assignExprs {
		assignmentQuery, assignmentArgs, errs := expr.AssignExpr()
		if len(errs) != 0 {
			return "", nil, errs[0]
		}

		assignments = append(assignments, assignmentQuery)
		args = append(args, assignmentArgs...)
	}
	sb.WriteString(strings.Join(assignments, ", "))

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
