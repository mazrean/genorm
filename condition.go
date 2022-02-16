package genorm

import (
	"fmt"
	"strings"
)

func And[T TableBase](expr1 TypedTableExpr[T, bool], expr2 TypedTableExpr[T, bool]) TypedTableExpr[T, bool] {
	errs1 := expr1.Errors()
	errs2 := expr2.Errors()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, bool]{
			errs: append(errs1, errs2...),
		}
	}

	query1, args1 := expr1.Expr()
	query2, args2 := expr2.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s AND %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func Or[T TableBase](expr1 TypedTableExpr[T, bool], expr2 TypedTableExpr[T, bool]) TypedTableExpr[T, bool] {
	errs1 := expr1.Errors()
	errs2 := expr2.Errors()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, bool]{
			errs: append(errs1, errs2...),
		}
	}

	query1, args1 := expr1.Expr()
	query2, args2 := expr2.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s OR %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func Xor[T TableBase](expr1 TypedTableExpr[T, bool], expr2 TypedTableExpr[T, bool]) TypedTableExpr[T, bool] {
	errs1 := expr1.Errors()
	errs2 := expr2.Errors()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, bool]{
			errs: append(errs1, errs2...),
		}
	}

	query1, args1 := expr1.Expr()
	query2, args2 := expr2.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s XOR %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func Not[T TableBase](expr TypedTableExpr[T, bool]) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(NOT %s)", query),
		args:  args,
	}
}

func Eq[T TableBase, S ExprType](expr1 TypedTableExpr[T, S], expr2 TypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs1 := expr1.Errors()
	errs2 := expr2.Errors()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, bool]{
			errs: append(errs1, errs2...),
		}
	}

	query1, args1 := expr1.Expr()
	query2, args2 := expr2.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s = %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func EqWithConst[T TableBase, S ExprType](expr TypedTableExpr[T, S], constant S) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s = ?)", query),
		args:  append(args, constant),
	}
}

func Neq[T TableBase, S ExprType](expr1 TypedTableExpr[T, S], expr2 TypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs1 := expr1.Errors()
	errs2 := expr2.Errors()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, bool]{
			errs: append(errs1, errs2...),
		}
	}

	query1, args1 := expr1.Expr()
	query2, args2 := expr2.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s != %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func NeqWithConst[T TableBase, S ExprType](expr TypedTableExpr[T, S], constant S) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s != ?)", query),
		args:  append(args, constant),
	}
}

func Leq[T TableBase, S ExprType](expr1 TypedTableExpr[T, S], expr2 TypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs1 := expr1.Errors()
	errs2 := expr2.Errors()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, bool]{
			errs: append(errs1, errs2...),
		}
	}

	query1, args1 := expr1.Expr()
	query2, args2 := expr2.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s <= %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func LeqWithConst[T TableBase, S ExprType](expr TypedTableExpr[T, S], constant S) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s <= ?)", query),
		args:  append(args, constant),
	}
}

func Geq[T TableBase, S ExprType](expr1 TypedTableExpr[T, S], expr2 TypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs1 := expr1.Errors()
	errs2 := expr2.Errors()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, bool]{
			errs: append(errs1, errs2...),
		}
	}

	query1, args1 := expr1.Expr()
	query2, args2 := expr2.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s >= %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func GeqWithConst[T TableBase, S ExprType](expr TypedTableExpr[T, S], constant S) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s >= ?)", query),
		args:  append(args, constant),
	}
}

func Lt[T TableBase, S ExprType](expr1 TypedTableExpr[T, S], expr2 TypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs1 := expr1.Errors()
	errs2 := expr2.Errors()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, bool]{
			errs: append(errs1, errs2...),
		}
	}

	query1, args1 := expr1.Expr()
	query2, args2 := expr2.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s < %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func LtWithConst[T TableBase, S ExprType](expr TypedTableExpr[T, S], constant S) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s < ?)", query),
		args:  append(args, constant),
	}
}

func Gt[T TableBase, S ExprType](expr1 TypedTableExpr[T, S], expr2 TypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs1 := expr1.Errors()
	errs2 := expr2.Errors()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, bool]{
			errs: append(errs1, errs2...),
		}
	}

	query1, args1 := expr1.Expr()
	query2, args2 := expr2.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s > %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func GtWithConst[T TableBase, S ExprType](expr TypedTableExpr[T, S], constant S) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s > ?)", query),
		args:  append(args, constant),
	}
}

func IsNull[T TableBase, S ExprType](expr NullableTypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s IS NULL)", query),
		args:  args,
	}
}

func IsNotNull[T TableBase, S ExprType](expr NullableTypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()
	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s IS NOT NULL)", query),
		args:  args,
	}
}

func In[T TableBase, S ExprType](expr1 TypedTableExpr[T, S], exprs ...TypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs := []error{}

	errs1 := expr1.Errors()
	if len(errs1) != 0 {
		errs = errs1
	}

	for _, expr := range exprs {
		errs2 := expr.Errors()
		if len(errs2) != 0 {
			errs = append(errs, errs2...)
		}
	}

	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query1, args1 := expr1.Expr()

	queries := []string{}
	args := []interface{}{}
	for _, expr := range exprs {
		query, args2 := expr.Expr()
		queries = append(queries, query)
		args = append(args, args2...)
	}

	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s IN (%s))", query1, strings.Join(queries, ", ")),
		args:  append(args1, args...),
	}
}

func InWithConst[T TableBase, S ExprType](expr TypedTableExpr[T, S], constants ...S) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()

	for _, constant := range constants {
		args = append(args, constant)
	}

	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s IN (%s))", query, strings.Repeat("?, ", len(constants)-1)+"?"),
		args:  args,
	}
}

func NotIn[T TableBase, S ExprType](expr1 TypedTableExpr[T, S], exprs ...TypedTableExpr[T, S]) TypedTableExpr[T, bool] {
	errs := []error{}

	errs1 := expr1.Errors()
	if len(errs1) != 0 {
		errs = errs1
	}

	for _, expr := range exprs {
		errs2 := expr.Errors()
		if len(errs2) != 0 {
			errs = append(errs, errs2...)
		}
	}

	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query1, args1 := expr1.Expr()

	queries := []string{}
	args := []interface{}{}
	for _, expr := range exprs {
		query, args2 := expr.Expr()
		queries = append(queries, query)
		args = append(args, args2...)
	}

	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s NOT IN (%s))", query1, strings.Join(queries, ", ")),
		args:  append(args1, args...),
	}
}

func NotInWithConst[T TableBase, S ExprType](expr TypedTableExpr[T, S], constants ...S) TypedTableExpr[T, bool] {
	errs := expr.Errors()
	if len(errs) != 0 {
		return &ExprStruct[T, bool]{
			errs: errs,
		}
	}

	query, args := expr.Expr()

	for _, constant := range constants {
		args = append(args, constant)
	}

	return &ExprStruct[T, bool]{
		query: fmt.Sprintf("(%s NOT IN (%s))", query, strings.Repeat("?, ", len(constants)-1)+"?"),
		args:  args,
	}
}
