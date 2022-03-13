package genorm

func NewWhereConditionClause[T Table](condition TypedTableExpr[T, WrappedPrimitive[bool]]) *whereConditionClause[T] {
	return &whereConditionClause[T]{
		condition: condition,
	}
}

func (c *whereConditionClause[T]) GetCondition() TypedTableExpr[T, WrappedPrimitive[bool]] {
	return c.condition
}

func (c *whereConditionClause[T]) Set(condition TypedTableExpr[T, WrappedPrimitive[bool]]) error {
	return c.set(condition)
}

func (c *whereConditionClause[T]) Exists() bool {
	return c.exists()
}

func (c *whereConditionClause[T]) GetExpr() (string, []ExprType, error) {
	return c.getExpr()
}

func NewGroupClause[T Table](exprs []TableExpr[T]) *groupClause[T] {
	return &groupClause[T]{
		exprs: exprs,
	}
}

func (c *groupClause[T]) GetCondition() []TableExpr[T] {
	return c.exprs
}

func (c *groupClause[T]) Set(exprs []TableExpr[T]) error {
	return c.set(exprs)
}

func (c *groupClause[T]) Exists() bool {
	return c.exists()
}

func (c *groupClause[T]) GetExpr() (string, []ExprType, error) {
	return c.getExpr()
}

type OrderItem[T Table] orderItem[T]

func NewOrderItem[T Table](expr TableExpr[T], direction OrderDirection) OrderItem[T] {
	return OrderItem[T]{
		expr:      expr,
		direction: direction,
	}
}

func (c *OrderItem[T]) Value() (TableExpr[T], OrderDirection) {
	return c.expr, c.direction
}

func NewOrderClause[T Table](items []OrderItem[T]) *orderClause[T] {
	orderItems := make([]orderItem[T], 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, orderItem[T](item))
	}

	return &orderClause[T]{
		orderExprs: orderItems,
	}
}

func (c *orderClause[T]) GetItems() []OrderItem[T] {
	items := make([]OrderItem[T], 0, len(c.orderExprs))
	for _, item := range c.orderExprs {
		items = append(items, OrderItem[T](item))
	}

	return items
}

func (c *orderClause[T]) Add(item OrderItem[T]) error {
	return c.add(orderItem[T](item))
}

func (c *orderClause[T]) Exists() bool {
	return c.exists()
}

func (c *orderClause[T]) GetExpr() (string, []ExprType, error) {
	return c.getExpr()
}

func NewLimitClause(limit uint64) *limitClause {
	return &limitClause{
		limit: limit,
	}
}

func (c *limitClause) GetLimit() uint64 {
	return c.limit
}

func (c *limitClause) Set(limit uint64) error {
	return c.set(limit)
}

func (c *limitClause) Exists() bool {
	return c.exists()
}

func (c *limitClause) GetExpr() (string, []ExprType, error) {
	return c.getExpr()
}
