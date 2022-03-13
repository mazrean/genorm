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
