package genorm

func (c *FindContext[_, _, _]) BuildQuery() (string, []ExprType, error) {
	return c.buildQuery()
}
