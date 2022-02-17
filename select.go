package genorm

import (
	"errors"
)

type SelectContext[T TableBase] struct {
	*Context[T]
	distinct        bool
	fields          []TableColumns[T]
	whereCondition  TypedTableExpr[T, bool]
	groupExpr       []TableExpr[T]
	havingCondition TypedTableExpr[T, bool]
	orderExprs      []orderItem[T]
	limit           uint64
	offset          uint64
	lockType        LockType
}

func NewSelectContext[T TableBase](table T) *SelectContext[T] {
	return &SelectContext[T]{
		Context: newContext(table),
	}
}

type orderItem[T TableBase] struct {
	expr      TableExpr[T]
	direction OrderDirection
}

type OrderDirection uint8

const (
	Asc OrderDirection = iota + 1
	Desc
)

type LockType uint8

const (
	none LockType = iota
	ForUpdate
	ForShare
)

func (c *SelectContext[TableBase]) Distinct() *SelectContext[TableBase] {
	if c.distinct {
		c.addError(errors.New("distinct already set"))
		return c
	}

	c.distinct = true

	return c
}

func (c *SelectContext[TableBase]) Fields(fields ...TableColumns[TableBase]) *SelectContext[TableBase] {
	if c.fields != nil {
		c.addError(errors.New("fields already set"))
		return c
	}
	if len(fields) == 0 {
		c.addError(errors.New("no fields"))
		return c
	}

	fields = append(c.fields, fields...)
	fieldMap := make(map[TableColumns[TableBase]]struct{}, len(fields))
	for _, field := range fields {
		if _, ok := fieldMap[field]; ok {
			c.addError(errors.New("duplicate field"))
			return c
		}

		fieldMap[field] = struct{}{}
	}

	c.fields = fields

	return c
}

func (c *SelectContext[TableBase]) Where(condition TypedTableExpr[TableBase, bool]) *SelectContext[TableBase] {
	if c.whereCondition != nil {
		c.addError(errors.New("where conditions already set"))
		return c
	}
	if condition == nil {
		c.addError(errors.New("empty where condition"))
		return c
	}

	c.whereCondition = condition

	return c
}

func (c *SelectContext[TableBase]) GroupBy(exprs ...TableExpr[TableBase]) *SelectContext[TableBase] {
	if len(exprs) == 0 {
		c.addError(errors.New("no group expr"))
		return c
	}

	c.groupExpr = append(c.groupExpr, exprs...)

	return c
}

func (c *SelectContext[TableBase]) Having(condition TypedTableExpr[TableBase, bool]) *SelectContext[TableBase] {
	if c.havingCondition != nil {
		c.addError(errors.New("having conditions already set"))
		return c
	}
	if condition == nil {
		c.addError(errors.New("empty having condition"))
		return c
	}

	c.havingCondition = condition

	return c
}

func (c *SelectContext[TableBase]) OrderBy(direction OrderDirection, expr TableExpr[TableBase]) *SelectContext[TableBase] {
	if expr == nil {
		c.addError(errors.New("empty order expr"))
		return c
	}

	if direction != Asc && direction != Desc {
		c.addError(errors.New("invalid order direction"))
		return c
	}

	c.orderExprs = append(c.orderExprs, orderItem[TableBase]{
		expr:      expr,
		direction: direction,
	})

	return c
}

func (c *SelectContext[TableBase]) Limit(limit uint64) *SelectContext[TableBase] {
	if c.limit != 0 {
		c.addError(errors.New("limit already set"))
		return c
	}
	if limit == 0 {
		c.addError(errors.New("invalid limit"))
		return c
	}

	c.limit = limit

	return c
}

func (c *SelectContext[TableBase]) Offset(offset uint64) *SelectContext[TableBase] {
	if c.offset != 0 {
		c.addError(errors.New("offset already set"))
		return c
	}

	c.offset = offset

	return c
}

func (c *SelectContext[TableBase]) Lock(lockType LockType) *SelectContext[TableBase] {
	if c.lockType != none {
		c.addError(errors.New("lock already set"))
		return c
	}

	if lockType != ForUpdate && lockType != ForShare {
		c.addError(errors.New("invalid lock type"))
		return c
	}

	c.lockType = lockType

	return c
}
