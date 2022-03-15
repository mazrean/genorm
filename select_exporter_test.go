package genorm

func (c *SelectContext[T]) BuildQuery() ([]Column, string, []ExprType, error) {
	return c.buildQuery()
}
