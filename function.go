package genorm

import (
	"errors"
	"fmt"
)

// Aggregate Functions

/*
Avg
if (distinct) {return AVG(DISTINCT expr)}
else {return AVG(expr)}
*/
func Avg[T Table, S ExprType](expr TypedTableExpr[T, S], distinct bool) TypedTableExpr[T, WrappedPrimitive[float64]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[float64]]{
			errs: []error{errors.New("avg expr is nil")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[float64]]{
			errs: errs,
		}
	}

	if distinct {
		query = fmt.Sprintf("AVG(DISTINCT %s)", query)
	} else {
		query = fmt.Sprintf("AVG(%s)", query)
	}

	return &ExprStruct[T, WrappedPrimitive[float64]]{
		query: query,
		args:  args,
	}
}

/*
Count
if (distinct) {return COUNT(DISTINCT expr)}
else {return COUNT(expr)}
*/
func Count[T Table, S ExprType](expr TypedTableExpr[T, S], distinct bool) TypedTableExpr[T, WrappedPrimitive[int64]] {
	if expr == nil {
		return &ExprStruct[T, WrappedPrimitive[int64]]{
			errs: []error{errors.New("count expr is nil")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, WrappedPrimitive[int64]]{
			errs: errs,
		}
	}

	if distinct {
		query = fmt.Sprintf("COUNT(DISTINCT %s)", query)
	} else {
		query = fmt.Sprintf("COUNT(%s)", query)
	}

	return &ExprStruct[T, WrappedPrimitive[int64]]{
		query: query,
		args:  args,
	}
}

// Max MAX(expr)
func Max[T Table, S ExprType](expr TypedTableExpr[T, S]) TypedTableExpr[T, S] {
	if expr == nil {
		return &ExprStruct[T, S]{
			errs: []error{errors.New("max expr is nil")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, S]{
			errs: errs,
		}
	}

	return &ExprStruct[T, S]{
		query: fmt.Sprintf("MAX(%s)", query),
		args:  args,
	}
}

// Min MIN(expr)
func Min[T Table, S ExprType](expr TypedTableExpr[T, S]) TypedTableExpr[T, S] {
	if expr == nil {
		return &ExprStruct[T, S]{
			errs: []error{errors.New("max expr is nil")},
		}
	}

	query, args, errs := expr.Expr()
	if len(errs) != 0 {
		return &ExprStruct[T, S]{
			errs: errs,
		}
	}

	return &ExprStruct[T, S]{
		query: fmt.Sprintf("MIN(%s)", query),
		args:  args,
	}
}
