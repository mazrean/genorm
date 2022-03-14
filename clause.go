package genorm

import (
	"errors"
	"fmt"
	"strings"
)

type whereConditionClause[T Table] struct {
	condition TypedTableExpr[T, WrappedPrimitive[bool]]
}

func (c *whereConditionClause[T]) set(condition TypedTableExpr[T, WrappedPrimitive[bool]]) error {
	if c.condition != nil {
		return errors.New("where conditions already set")
	}
	if condition == nil {
		return errors.New("empty where condition")
	}

	c.condition = condition

	return nil
}

func (c *whereConditionClause[T]) exists() bool {
	return c.condition != nil
}

func (c *whereConditionClause[T]) getExpr() (string, []ExprType, error) {
	if c.condition == nil {
		return "", nil, errors.New("empty where condition")
	}

	query, args, errs := c.condition.Expr()
	if len(errs) != 0 {
		return "", nil, errs[0]
	}

	return query, args, nil
}

type groupClause[T Table] struct {
	exprs []TableExpr[T]
}

func (c *groupClause[T]) set(exprs []TableExpr[T]) error {
	if len(c.exprs) != 0 {
		return errors.New("group by already set")
	}
	if len(exprs) == 0 {
		return errors.New("empty group by")
	}

	c.exprs = exprs

	return nil
}

func (c *groupClause[_]) exists() bool {
	return len(c.exprs) != 0
}

func (c *groupClause[T]) getExpr() (string, []ExprType, error) {
	if len(c.exprs) == 0 {
		return "", nil, errors.New("empty group by")
	}

	queries := make([]string, 0, len(c.exprs))
	args := []ExprType{}
	for _, expr := range c.exprs {
		groupQuery, groupArgs, errs := expr.Expr()
		if len(errs) != 0 {
			return "", nil, errs[0]
		}

		queries = append(queries, groupQuery)
		args = append(args, groupArgs...)
	}

	return "GROUP BY " + strings.Join(queries, ", "), args, nil
}

type orderClause[T Table] struct {
	orderExprs []orderItem[T]
}

type orderItem[T Table] struct {
	expr      TableExpr[T]
	direction OrderDirection
}

type OrderDirection uint8

const (
	Asc OrderDirection = iota + 1
	Desc
)

func (c *orderClause[T]) add(item orderItem[T]) error {
	if item.expr == nil {
		return errors.New("empty order expr")
	}

	if item.direction != Asc && item.direction != Desc {
		return errors.New("invalid order direction")
	}

	c.orderExprs = append(c.orderExprs, item)

	return nil
}

func (c *orderClause[T]) exists() bool {
	return len(c.orderExprs) != 0
}

func (c *orderClause[T]) getExpr() (string, []ExprType, error) {
	if len(c.orderExprs) == 0 {
		return "", nil, errors.New("empty order by")
	}

	args := []ExprType{}
	orderQueries := make([]string, 0, len(c.orderExprs))
	for _, orderItem := range c.orderExprs {
		orderQuery, orderArgs, errs := orderItem.expr.Expr()
		if len(errs) != 0 {
			return "", nil, errs[0]
		}

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

	return fmt.Sprintf("ORDER BY %s", strings.Join(orderQueries, ", ")), args, nil
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

func (l *limitClause) getExpr() (string, []ExprType, error) {
	if l.limit == 0 {
		return "", nil, errors.New("empty limit")
	}

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

func (o *offsetClause) getExpr() (string, []ExprType, error) {
	if o.offset == 0 {
		return "", nil, errors.New("empty offset")
	}

	return fmt.Sprintf("OFFSET %d", o.offset), nil, nil
}

type lockClause struct {
	lockType LockType
}

type LockType uint8

const (
	none LockType = iota
	ForUpdate
	ForShare
)

func (l *lockClause) set(lockType LockType) error {
	if l.lockType != none {
		return errors.New("lock type already set")
	}
	if lockType != ForUpdate && lockType != ForShare {
		return errors.New("invalid lock type")
	}

	l.lockType = lockType

	return nil
}

func (l *lockClause) exists() bool {
	return l.lockType != none
}

func (l *lockClause) getExpr() (string, []ExprType, error) {
	switch l.lockType {
	case ForUpdate:
		return "FOR UPDATE", nil, nil
	case ForShare:
		return "FOR SHARE", nil, nil
	}

	return "", nil, errors.New("invalid lock type")
}
