package genorm

func (c *InsertContext[T]) BuildQuery() (string, []any, error) {
	return c.buildQuery()
}
