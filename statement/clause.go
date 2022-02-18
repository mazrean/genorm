package statement

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mazrean/genorm"
)

type whereConditionClause[T Table] struct {
	condition genorm.TypedTableExpr[T, bool]
}

func (c *whereConditionClause[Table]) set(condition genorm.TypedTableExpr[Table, bool]) error {
	if c.condition != nil {
		return errors.New("where conditions already set")
	}
	if condition == nil {
		return errors.New("empty where condition")
	}

	c.condition = condition

	return nil
}

func (c *whereConditionClause[Table]) exists() bool {
	return c.condition != nil
}

func (c *whereConditionClause[Table]) getExpr() (string, []any, error) {
	query, args := c.condition.Expr()

	return query, args, nil
}

type orderClause[T Table] struct {
	orderExprs []orderItem[T]
}

type orderItem[T Table] struct {
	expr      genorm.TableExpr[T]
	direction OrderDirection
}

type OrderDirection uint8

const (
	Asc OrderDirection = iota + 1
	Desc
)

func (c *orderClause[Table]) add(item orderItem[Table]) error {
	if item.expr == nil {
		return errors.New("empty order expr")
	}

	if item.direction != Asc && item.direction != Desc {
		return errors.New("invalid order direction")
	}

	c.orderExprs = append(c.orderExprs, item)

	return nil
}

func (c *orderClause[Table]) exists() bool {
	return len(c.orderExprs) != 0
}

func (c *orderClause[Table]) getExpr() (string, []any, error) {
	args := []any{}
	orderQueries := make([]string, 0, len(c.orderExprs))
	for _, orderItem := range c.orderExprs {
		orderQuery, orderArgs := orderItem.expr.Expr()

		var directionQuery string
		switch orderItem.direction {
		case Asc:
			directionQuery = "ASC"
		case Desc:
			directionQuery = "DESC"
		default:
			return "", nil, fmt.Errorf("invalid order direction: %d", orderItem.direction)
		}

		orderQueries = append(orderQueries, fmt.Sprintf("%s %s", orderQuery, directionQuery))
		args = append(args, orderArgs...)
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(orderQueries, ", ")), nil, nil
}

type limitClause struct {
	limit uint64
}

func (l *limitClause) set(limit uint64) error {
	if l.limit != 0 {
		return errors.New("limit already set")
	}
	if limit == 0 {
		return errors.New("invalid limit")
	}

	l.limit = limit

	return nil
}

func (l *limitClause) exists() bool {
	return l.limit != 0
}

func (l *limitClause) getExpr() (string, []any, error) {
	return fmt.Sprintf("LIMIT %d", l.limit), nil, nil
}

type offsetClause struct {
	offset uint64
}

func (o *offsetClause) set(offset uint64) error {
	if o.offset != 0 {
		return errors.New("offset already set")
	}
	if offset == 0 {
		return errors.New("invalid offset")
	}

	o.offset = offset

	return nil
}

func (o *offsetClause) exists() bool {
	return o.offset != 0
}

func (o *offsetClause) getExpr() (string, []any, error) {
	return fmt.Sprintf("OFFSET %d", o.offset), nil, nil
}
