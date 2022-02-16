package genorm

type ExprStruct[T TableBase, S ExprType] struct {
	query string
	args  []any
	errs  []error
}

func RawExpr[T TableBase, S ExprType](query string, args ...any) *ExprStruct[T, S] {
	return &ExprStruct[T, S]{
		query: query,
		args:  args,
	}
}

func (es *ExprStruct[TableBase, ExprType]) Expr() (string, []any) {
	return es.query, es.args
}

func (es *ExprStruct[TableBase, ExprType]) Errors() []error {
	return es.errs
}

func (es *ExprStruct[TableBase, ExprType]) TableExpr(TableBase) (string, []any) {
	return es.Expr()
}

func (es *ExprStruct[TableBase, ExprType]) TypedExpr(ExprType) (string, []any) {
	return es.Expr()
}
