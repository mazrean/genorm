package genorm

import (
	"errors"
	"fmt"
	"strings"
)

// Assign Operators

// Assign expr1 = expr2
func Assign[T Table, S ExprType](
	expr1 TypedTableColumns[T, S],
	expr2 TypedTableExpr[T, S],
) *TableAssignExpr[T] {
	if expr1 == nil || expr2 == nil {
		return &TableAssignExpr[T]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &TableAssignExpr[T]{
			errs: append(errs1, errs2...),
		}
	}

	return &TableAssignExpr[T]{
		query: fmt.Sprintf("%s = %s", query1, query2),
		args:  append(args1, args2...),
	}
}

// AssignLit expr = literal
func AssignLit[T Table, S ExprType](
	expr TypedTableColumns[T, S],
	literal S,
) *TableAssignExpr[T] {
	if expr == nil {
		return &TableAssignExpr[T]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &TableAssignExpr[T]{
			errs: errs,
		}
	}

	return &TableAssignExpr[T]{
		query: fmt.Sprintf("%s = ?", query),
		args:  append(args, literal),
	}
}

// Logical Operators

// And (expr1 AND expr2)
func And[T Table](
	expr1 TypedTableExpr[T, WrappedPrimitive[bool]],
	expr2 TypedTableExpr[T, WrappedPrimitive[bool]],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr1 == nil || expr2 == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s AND %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

// Or (expr1 OR expr2)
func Or[T Table](
	expr1 TypedTableExpr[T, WrappedPrimitive[bool]],
	expr2 TypedTableExpr[T, WrappedPrimitive[bool]],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr1 == nil || expr2 == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s OR %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

// Xor (expr1 XOR expr2)
func Xor[T Table](
	expr1 TypedTableExpr[T, WrappedPrimitive[bool]],
	expr2 TypedTableExpr[T, WrappedPrimitive[bool]],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr1 == nil || expr2 == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s XOR %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

// Not (NOT expr)
func Not[T Table](
	expr TypedTableExpr[T, WrappedPrimitive[bool]],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(NOT %s)", query),
		args:  args,
	}
}

// Comparison Operators

// Eq (expr1 = expr2)
func Eq[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr1 == nil || expr2 == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s = %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

// EqLit (expr = literal)
func EqLit[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	literal S,
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s = ?)", query),
		args:  append(args, literal),
	}
}

// Neq (expr1 != expr2)
func Neq[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr1 == nil || expr2 == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s != %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

// NeqLit (expr != literal)
func NeqLit[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	literal S,
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s != ?)", query),
		args:  append(args, literal),
	}
}

// Leq (expr1 <= expr2)
func Leq[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr1 == nil || expr2 == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s <= %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

// LeqLit (expr <= literal)
func LeqLit[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	literal S,
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s <= ?)", query),
		args:  append(args, literal),
	}
}

// Geq (expr1 >= expr2)
func Geq[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr1 == nil || expr2 == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s >= %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

// GeqLit (expr >= literal)
func GeqLit[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	literal S,
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s >= ?)", query),
		args:  append(args, literal),
	}
}

// Lt (expr1 < expr2)
func Lt[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr1 == nil || expr2 == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s < %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

// LtLit (expr < literal)
func LtLit[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	literal S,
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s < ?)", query),
		args:  append(args, literal),
	}
}

// Gt (expr1 > expr2)
func Gt[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	expr2 TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr1 == nil || expr2 == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query1, args1, errs1 := expr1.Expr()
	query2, args2, errs2 := expr2.Expr()
	if len(errs1) != 0 || len(errs2) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: append(errs1, errs2...),
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s > %s)", query1, query2),
		args:  append(args1, args2...),
	}
}

// GtLit (expr > literal)
func GtLit[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	literal S,
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s > ?)", query),
		args:  append(args, literal),
	}
}

// IsNull (expr IS NULL)
func IsNull[T Table, S ExprType](
	expr TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s IS NULL)", query),
		args:  args,
	}
}

// IsNotNull (expr IS NOT NULL)
func IsNotNull[T Table, S ExprType](
	expr TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s IS NOT NULL)", query),
		args:  args,
	}
}

// In (expr1 IN (exprs[0], exprs[1], ...))
func In[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	exprs ...TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	errs := []error{}

	if expr1 == nil || exprs == nil {
		errs = append(errs, errors.New("Assign: nil expression"))
	}

	if len(exprs) == 0 {
		errs = append(errs, errors.New("invalid number of arguments for in"))
	}

	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	query1, args1, errs1 := expr1.Expr()
	if len(errs1) != 0 {
		errs = errs1
	}

	queries := []string{}
	args := []ExprType{}
	for _, expr := range exprs {
		query, args2, errs2 := expr.Expr()
		if len(errs2) != 0 {
			errs = append(errs, errs2...)
		}

		queries = append(queries, query)
		args = append(args, args2...)
	}

	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s IN (%s))", query1, strings.Join(queries, ", ")),
		args:  append(args1, args...),
	}
}

// InLit (expr IN (literals[0], literal[1], ...))
func InLit[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	literals ...S,
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil || literals == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	if len(literals) == 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("invalid number of arguments for in")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	for _, literal := range literals {
		args = append(args, literal)
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s IN (%s))", query, strings.Repeat("?, ", len(literals)-1)+"?"),
		args:  args,
	}
}

// NotIn (expr1 NOT IN (exprs[0], exprs[1], ...))
func NotIn[T Table, S ExprType](
	expr1 TypedTableExpr[T, S],
	exprs ...TypedTableExpr[T, S],
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	errs := []error{}

	if expr1 == nil || exprs == nil {
		errs = append(errs, errors.New("Assign: nil expression"))
	}

	if len(exprs) == 0 {
		errs = append(errs, errors.New("invalid number of arguments for in"))
	}

	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	query1, args1, errs1 := expr1.Expr()
	if len(errs1) != 0 {
		errs = errs1
	}

	queries := []string{}
	args := []ExprType{}
	for _, expr := range exprs {
		query, args2, errs2 := expr.Expr()
		if len(errs2) != 0 {
			errs = append(errs, errs2...)
		}

		queries = append(queries, query)
		args = append(args, args2...)
	}

	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s NOT IN (%s))", query1, strings.Join(queries, ", ")),
		args:  append(args1, args...),
	}
}

// NotInLit (expr NOT IN (literals[0], literal[1], ...))
func NotInLit[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	literals ...S,
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	if expr == nil || literals == nil {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("Assign: nil expression")},
		}
	}

	if len(literals) == 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: []error{errors.New("invalid number of arguments for in")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[bool]]{
			errs: errs,
		}
	}

	for _, literal := range literals {
		args = append(args, literal)
	}

	return &ExprStruct[T, WrappedPrimitive[bool]]{
		query: fmt.Sprintf("(%s NOT IN (%s))", query, strings.Repeat("?, ", len(literals)-1)+"?"),
		args:  args,
	}
}
