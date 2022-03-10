package genorm

type Expr interface {
	Expr() (string, []ExprType, []error)
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

type TableAssignExpr[T Table] struct {
	query string
	args  []ExprType
	errs  []error
}

func (tae *TableAssignExpr[_]) AssignExpr() (string, []ExprType, []error) {
	if len(tae.errs) != 0 {
		return "", nil, tae.errs
	}

	return tae.query, tae.args, nil
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

func (es *ExprStruct[T, _]) TableExpr(T) (string, []ExprType, []error) {
	return es.Expr()
}

func (es *ExprStruct[_, S]) TypedExpr(S) (string, []ExprType, []error) {
	return es.Expr()
}
