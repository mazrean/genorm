package genorm

import (
	"fmt"
	"strings"
)

// Logical Operators

func And[T Table](
	expr1 TypedTableExpr[T, *WrappedPrimitive[bool]],
	expr2 TypedTableExpr[T, *WrappedPrimitive[bool]],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s AND %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func Or[T Table](
	expr1 TypedTableExpr[T, *WrappedPrimitive[bool]],
	expr2 TypedTableExpr[T, *WrappedPrimitive[bool]],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s OR %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func Xor[T Table](
	expr1 TypedTableExpr[T, *WrappedPrimitive[bool]],
	expr2 TypedTableExpr[T, *WrappedPrimitive[bool]],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s XOR %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func Not[T Table](
	expr TypedTableExpr[T, *WrappedPrimitive[bool]],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(NOT %s)", query),
		args:  args,
	}
}

// Comparison Operators

func Eq[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s = %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func EqWithConst[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	constant S,
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s = ?)", query),
		args:  append(args, constant),
	}
}

func Neq[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s != %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func NeqWithConst[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	constant S,
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s != ?)", query),
		args:  append(args, constant),
	}
}

func Leq[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s <= %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func LeqWithConst[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	constant S,
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s <= ?)", query),
		args:  append(args, constant),
	}
}

func Geq[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s >= %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func GeqWithConst[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	constant S,
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s >= ?)", query),
		args:  append(args, constant),
	}
}

func Lt[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s < %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func LtWithConst[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	constant S,
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s < ?)", query),
		args:  append(args, constant),
	}
}

func Gt[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s > %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

func GtWithConst[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	constant S,
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s > ?)", query),
		args:  append(args, constant),
	}
}

func IsNull[T Table, S ExprType](
	expr NullableTypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s IS NULL)", query),
		args:  args,
	}
}

func IsNotNull[T Table, S ExprType](
	expr NullableTypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s IS NOT NULL)", query),
		args:  args,
	}
}

func In[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	exprs ...TypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	errs := []error{}

	query1, args1, errs1 := expr1.Expr()
	if len(errs1) != 0 {
		errs = errs1
	}

	queries := []string{}
	args := []interface{}{}
	for _, expr := range exprs {
		query, args2, errs2 := expr.Expr()
		if len(errs2) != 0 {
			errs = append(errs, errs2...)
		}

		queries = append(queries, query)
		args = append(args, args2...)
	}

	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s IN (%s))", query1, strings.Join(queries, ", ")),
		args:  append(args1, args...),
	}
}

func InWithConst[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	constants ...S,
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	for _, constant := range constants {
		args = append(args, constant)
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s IN (%s))", query, strings.Repeat("?, ", len(constants)-1)+"?"),
		args:  args,
	}
}

func NotIn[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	exprs ...TypedTableExpr[T, S],
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	errs := []error{}

	query1, args1, errs1 := expr1.Expr()
	if len(errs1) != 0 {
		errs = errs1
	}

	queries := []string{}
	args := []interface{}{}
	for _, expr := range exprs {
		query, args2, errs2 := expr.Expr()
		if len(errs2) != 0 {
			errs = append(errs, errs2...)
		}

		queries = append(queries, query)
		args = append(args, args2...)
	}

	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s NOT IN (%s))", query1, strings.Join(queries, ", ")),
		args:  append(args1, args...),
	}
}

func NotInWithConst[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	constants ...S,
) TypedTableExpr[T, *WrappedPrimitive[bool]] {
	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, *WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	for _, constant := range constants {
		args = append(args, constant)
	}

	return &ExprStruct[T, *WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s NOT IN (%s))", query, strings.Repeat("?, ", len(constants)-1)+"?"),
		args:  args,
	}
}