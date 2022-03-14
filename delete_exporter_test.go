package genorm

func (c *DeleteContext[T]) BuildQuery() (string, []ExprType, error) {
	return c.buildQuery()
}
