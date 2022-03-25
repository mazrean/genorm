package genorm

type Tuple interface {
	Exprs() []Expr
	Columns() []ColumnFieldExprType
}

type TuplePointer[T any] interface {
	Exprs() []Expr
	Columns() []ColumnFieldExprType
	*T
}

type Tuple2Struct[
	S Table,
	T1 ExprType, U1 ColumnFieldExprTypePointer[T1],
	T2 ExprType, U2 ColumnFieldExprTypePointer[T2],
] struct {
	value1 T1
	value2 T2
	expr1  TypedTableExpr[S, T1]
	expr2  TypedTableExpr[S, T2]
}

func Tuple2[
	S Table,
	T1 ExprType, U1 ColumnFieldExprTypePointer[T1],
	T2 ExprType, U2 ColumnFieldExprTypePointer[T2],
](expr1 TypedTableExpr[S, T1], expr2 TypedTableExpr[S, T2]) *Tuple2Struct[S, T1, U1, T2, U2] {
	return &Tuple2Struct[S, T1, U1, T2, U2]{
		expr1: expr1,
		expr2: expr2,
	}
}

func (t *Tuple2Struct[_, _, _, _, _]) Exprs() []Expr {
	return []Expr{t.expr1, t.expr2}
}

func (t *Tuple2Struct[_, _, U1, _, U2]) Columns() []ColumnFieldExprType {
	return []ColumnFieldExprType{
		U1(&t.value1),
		U2(&t.value2),
	}
}

func (t *Tuple2Struct[_, T1, _, T2, _]) Values() (T1, T2) {
	return t.value1, t.value2
}

type Tuple3Struct[
	S Table,
	T1 ExprType, U1 ColumnFieldExprTypePointer[T1],
	T2 ExprType, U2 ColumnFieldExprTypePointer[T2],
	T3 ExprType, U3 ColumnFieldExprTypePointer[T3],
] struct {
	value1 T1
	value2 T2
	value3 T3
	expr1  TypedTableExpr[S, T1]
	expr2  TypedTableExpr[S, T2]
	expr3  TypedTableExpr[S, T3]
}

func Tuple3[
	S Table,
	T1 ExprType, U1 ColumnFieldExprTypePointer[T1],
	T2 ExprType, U2 ColumnFieldExprTypePointer[T2],
	T3 ExprType, U3 ColumnFieldExprTypePointer[T3],
](expr1 TypedTableExpr[S, T1], expr2 TypedTableExpr[S, T2], expr3 TypedTableExpr[S, T3]) *Tuple3Struct[S, T1, U1, T2, U2, T3, U3] {
	return &Tuple3Struct[S, T1, U1, T2, U2, T3, U3]{
		expr1: expr1,
		expr2: expr2,
		expr3: expr3,
	}
}

func (t *Tuple3Struct[_, _, _, _, _, _, _]) Exprs() []Expr {
	return []Expr{t.expr1, t.expr2, t.expr3}
}

func (t *Tuple3Struct[_, _, U1, _, U2, _, U3]) Columns() []ColumnFieldExprType {
	return []ColumnFieldExprType{
		U1(&t.value1),
		U2(&t.value2),
		U3(&t.value3),
	}
}

func (t *Tuple3Struct[_, T1, _, T2, _, T3, _]) Values() (T1, T2, T3) {
	return t.value1, t.value2, t.value3
}

type Tuple4Struct[
	S Table,
	T1 ExprType, U1 ColumnFieldExprTypePointer[T1],
	T2 ExprType, U2 ColumnFieldExprTypePointer[T2],
	T3 ExprType, U3 ColumnFieldExprTypePointer[T3],
	T4 ExprType, U4 ColumnFieldExprTypePointer[T4],
] struct {
	value1 T1
	value2 T2
	value3 T3
	value4 T4
	expr1  TypedTableExpr[S, T1]
	expr2  TypedTableExpr[S, T2]
	expr3  TypedTableExpr[S, T3]
	expr4  TypedTableExpr[S, T4]
}

func Tuple4[
	S Table,
	T1 ExprType, U1 ColumnFieldExprTypePointer[T1],
	T2 ExprType, U2 ColumnFieldExprTypePointer[T2],
	T3 ExprType, U3 ColumnFieldExprTypePointer[T3],
	T4 ExprType, U4 ColumnFieldExprTypePointer[T4],
](expr1 TypedTableExpr[S, T1], expr2 TypedTableExpr[S, T2], expr3 TypedTableExpr[S, T3], expr4 TypedTableExpr[S, T4]) *Tuple4Struct[S, T1, U1, T2, U2, T3, U3, T4, U4] {
	return &Tuple4Struct[S, T1, U1, T2, U2, T3, U3, T4, U4]{
		expr1: expr1,
		expr2: expr2,
		expr3: expr3,
		expr4: expr4,
	}
}

func (t *Tuple4Struct[_, _, _, _, _, _, _, _, _]) Exprs() []Expr {
	return []Expr{t.expr1, t.expr2, t.expr3, t.expr4}
}

func (t *Tuple4Struct[_, _, U1, _, U2, _, U3, _, U4]) Columns() []ColumnFieldExprType {
	return []ColumnFieldExprType{
		U1(&t.value1),
		U2(&t.value2),
		U3(&t.value3),
		U4(&t.value4),
	}
}

func (t *Tuple4Struct[_, T1, _, T2, _, T3, _, T4, _]) Values() (T1, T2, T3, T4) {
	return t.value1, t.value2, t.value3, t.value4
}

type Tuple5Struct[
	S Table,
	T1 ExprType, U1 ColumnFieldExprTypePointer[T1],
	T2 ExprType, U2 ColumnFieldExprTypePointer[T2],
	T3 ExprType, U3 ColumnFieldExprTypePointer[T3],
	T4 ExprType, U4 ColumnFieldExprTypePointer[T4],
	T5 ExprType, U5 ColumnFieldExprTypePointer[T5],
] struct {
	value1 T1
	value2 T2
	value3 T3
	value4 T4
	value5 T5
	expr1  TypedTableExpr[S, T1]
	expr2  TypedTableExpr[S, T2]
	expr3  TypedTableExpr[S, T3]
	expr4  TypedTableExpr[S, T4]
	expr5  TypedTableExpr[S, T5]
}

func Tuple5[
	S Table,
	T1 ExprType, U1 ColumnFieldExprTypePointer[T1],
	T2 ExprType, U2 ColumnFieldExprTypePointer[T2],
	T3 ExprType, U3 ColumnFieldExprTypePointer[T3],
	T4 ExprType, U4 ColumnFieldExprTypePointer[T4],
	T5 ExprType, U5 ColumnFieldExprTypePointer[T5],
](expr1 TypedTableExpr[S, T1], expr2 TypedTableExpr[S, T2], expr3 TypedTableExpr[S, T3], expr4 TypedTableExpr[S, T4], expr5 TypedTableExpr[S, T5]) *Tuple5Struct[S, T1, U1, T2, U2, T3, U3, T4, U4, T5, U5] {
	return &Tuple5Struct[S, T1, U1, T2, U2, T3, U3, T4, U4, T5, U5]{
		expr1: expr1,
		expr2: expr2,
		expr3: expr3,
		expr4: expr4,
		expr5: expr5,
	}
}

func (t *Tuple5Struct[_, _, _, _, _, _, _, _, _, _, _]) Exprs() []Expr {
	return []Expr{t.expr1, t.expr2, t.expr3, t.expr4, t.expr5}
}

func (t *Tuple5Struct[_, _, U1, _, U2, _, U3, _, U4, _, U5]) Columns() []ColumnFieldExprType {
	return []ColumnFieldExprType{
		U1(&t.value1),
		U2(&t.value2),
		U3(&t.value3),
		U4(&t.value4),
		U5(&t.value5),
	}
}

func (t *Tuple5Struct[_, T1, _, T2, _, T3, _, T4, _, T5, _]) Values() (T1, T2, T3, T4, T5) {
	return t.value1, t.value2, t.value3, t.value4, t.value5
}
