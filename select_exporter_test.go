package genorm

func (c *SelectContext[_, _]) BuildQuery() ([]Column, string, []ExprType, error) {
	return c.buildQuery()
}
