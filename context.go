package genorm

type Context[T Table] struct {
	table  T
	errs   []error
	config genormConfig
}

func newContext[T Table](table T, options ...Option) *Context[T] {
	config := defaultConfig()
	for _, option := range options {
		option(&config)
	}

	return &Context[T]{
		table:  table,
		errs:   table.GetErrors(),
		config: config,
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
