package genorm

func NewTableAssignExpr[T Table](
	query string,
	args []ExprType,
	errs []error,
) *TableAssignExpr[T] {
	return &TableAssignExpr[T]{
		query: query,
		args:  args,
		errs:  errs,
	}
}
