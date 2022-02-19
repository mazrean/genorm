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
	assignmentMap  map[genorm.TableColumns[T]]genorm.TableExpr[T]
	whereCondition whereConditionClause[T]
	order          orderClause[T]
	limit          limitClause
}

func NewUpdateContext[T Table](table T) *UpdateContext[T] {
	ctx := newContext(table)

	return &UpdateContext[T]{
		Context:       ctx,
		assignmentMap: map[genorm.TableColumns[T]]genorm.TableExpr[T]{},
	}
}

func (c *UpdateContext[Table]) Assign(column genorm.TableColumns[Table], expr genorm.TableExpr[Table]) (res *UpdateContext[Table]) {
	res = c

	defer func() {
		err := recover()
		if err != nil {
			c.addError(errors.New("column should be comparable"))
		}
	}()

	if column == nil {
		c.addError(errors.New("invalid column"))
		return c
	}
	if expr == nil {
		c.addError(errors.New("invalid expression"))
		return c
	}
	if c.assignmentMap == nil {
		c.addError(errors.New("assignment map not initialized"))
	}

	if _, ok := c.assignmentMap[column]; ok { // if column is not comparable, panic
		c.addError(errors.New("duplicate assignment"))
	}

	c.assignmentMap[column] = expr

	return c
}

func (c *UpdateContext[Table]) MapAssign(assignMap map[genorm.TableColumns[Table]]genorm.TableExpr[Table]) *UpdateContext[Table] {
	for column, expr := range assignMap {
		if column == nil {
			c.addError(errors.New("invalid column"))
			return c
		}
		if expr == nil {
			c.addError(errors.New("invalid expression"))
			return c
		}
		if c.assignmentMap == nil {
			c.addError(errors.New("assignment map not initialized"))
		}

		if _, ok := c.assignmentMap[column]; ok {
			c.addError(errors.New("duplicate assignment"))
		}

		c.assignmentMap[column] = expr
	}

	return c
}

func (c *UpdateContext[Table]) Where(
	condition genorm.TypedTableExpr[Table, *genorm.WrappedPrimitive[bool]],
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

func (c *UpdateContext[Table]) Do(db *sql.DB) (rowsAffected int64, err error) {
	return c.DoContext(context.Background(), db)
}

func (c *UpdateContext[Table]) buildQuery() (string, []any, error) {
	args := []any{}

	sb := strings.Builder{}
	sb.WriteString("UPDATE ")
	sb.WriteString(c.table.SQLTableName())

	if len(c.assignmentMap) == 0 {
		return "", nil, errors.New("no assignment")
	}

	sb.WriteString(" SET ")
	assignments := make([]string, 0, len(c.assignmentMap))
	for column, expr := range c.assignmentMap {
		assignmentQuery, assignmentArgs := expr.Expr()
		assignments = append(assignments, fmt.Sprintf("%s = %s", column.SQLColumnName(), assignmentQuery))
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
