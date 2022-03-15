package genorm

func (c *PluckContext[T, _]) BuildQuery() (string, []ExprType, error) {
	return c.buildQuery()
}
