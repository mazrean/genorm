package genorm

type Expr interface {
	Expr() (string, []any)
	Errors() []error
}

type NullableExpr interface {
	Expr
	IsNull() bool
}

type TableExpr[T Table] interface {
	Expr
	TableExpr(T) (string, []any)
}

type TypedExpr[T ExprType] interface {
	Expr
	TypedExpr(T) (string, []any)
}

type TypedTableExpr[T Table, S ExprType] interface {
	Expr
	TableExpr[T]
	TypedExpr[S]
}

type NullableTypedTableExpr[T Table, S ExprType] interface {
	NullableExpr
	TableExpr[T]
	TypedExpr[S]
}

type ExprStruct[T Table, S ExprType] struct {
	query string
	args  []any
	errs  []error
}

func RawExpr[T Table, S ExprType](query string, args ...any) *ExprStruct[T, S] {
	return &ExprStruct[T, S]{
		query: query,
		args:  args,
	}
}

func (es *ExprStruct[_, _]) Expr() (string, []any) {
	return es.query, es.args
}

func (es *ExprStruct[_, _]) Errors() []error {
	return es.errs
}

func (es *ExprStruct[Table, _]) TableExpr(Table) (string, []any) {
	return es.Expr()
}

func (es *ExprStruct[_, ExprType]) TypedExpr(ExprType) (string, []any) {
	return es.Expr()
}
