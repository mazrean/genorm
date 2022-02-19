package genorm

type Expr interface {
	Expr() (string, []ExprType, []error)
}

type NullableExpr interface {
	Expr
	IsNull() bool
}

type TableExpr[T Table] interface {
	Expr
	TableExpr(T) (string, []ExprType, []error)
}

type TypedExpr[T ExprType] interface {
	Expr
	TypedExpr(T) (string, []ExprType, []error)
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
	args  []ExprType
	errs  []error
}

func RawExpr[T Table, S ExprType](query string, args ...ExprType) *ExprStruct[T, S] {
	return &ExprStruct[T, S]{
		query: query,
		args:  args,
	}
}

func (es *ExprStruct[_, _]) Expr() (string, []ExprType, []error) {
	if len(es.errs) != 0 {
		return "", nil, es.errs
	}

	return es.query, es.args, nil
}

func (es *ExprStruct[Table, _]) TableExpr(Table) (string, []ExprType, []error) {
	return es.Expr()
}

func (es *ExprStruct[_, Type]) TypedExpr(Type) (string, []ExprType, []error) {
	return es.Expr()
}
