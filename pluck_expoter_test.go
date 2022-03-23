package genorm

func (c *PluckContext[_, _]) BuildQuery() (string, []ExprType, error) {
	return c.buildQuery()
}
