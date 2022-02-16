package genorm

type Context[T TableBase] struct {
	table T
	errs  []error
}

func newContext[T TableBase](table T) *Context[T] {
	return &Context[T]{
		table: table,
	}
}

func (c *Context[T]) Table() T {
	return c.table
}

func (c *Context[T]) addError(err error) {
	c.errs = append(c.errs, err)
}

func (c *Context[T]) Errors() []error {
	if len(c.errs) == 0 {
		return nil
	}

	return c.errs
}
