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
