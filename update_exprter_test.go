package genorm

func (c *UpdateContext[T]) BuildQuery() (string, []ExprType, error) {
	return c.buildQuery()
}
