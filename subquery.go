package genorm

import (
	"fmt"
)

// SubQueryExpr represents a subquery that returns a single value
// It can be used in WHERE clauses, HAVING clauses, or as a value in comparisons
type SubQueryExpr[TOuter Table, TInner Table, S ExprType] struct {
	query string
	args  []ExprType
	errs  []error
}

// SubQuery creates a subquery expression from a PluckContext
// The subquery returns a single value and can be used in WHERE clauses or comparisons
// Example:
//
//	subQuery := genorm.SubQuery(genorm.Pluck(orm.InnerTable(), innerColumn).Where(...))
//	results, err := genorm.Select(orm.OuterTable()).
//	  Where(genorm.Eq(outerColumn, subQuery)).
//	  GetAll(db)
func SubQuery[TOuter Table, TInner Table, S ExprType](
	pluckCtx *PluckContext[TInner, S],
) TypedTableExpr[TOuter, S] {
	if pluckCtx == nil {
		return &SubQueryExpr[TOuter, TInner, S]{
			errs: []error{fmt.Errorf("subquery pluck context is nil")},
		}
	}

	errs := pluckCtx.Errors()
	if len(errs) != 0 {
		return &SubQueryExpr[TOuter, TInner, S]{
			errs: errs,
		}
	}

	query, args, err := pluckCtx.buildQuery()
	if err != nil {
		return &SubQueryExpr[TOuter, TInner, S]{
			errs: []error{err},
		}
	}

	// Wrap the query in parentheses to make it a valid subquery
	return &SubQueryExpr[TOuter, TInner, S]{
		query: fmt.Sprintf("(%s)", query),
		args:  args,
	}
}

// Expr implements the Expr interface
func (sq *SubQueryExpr[_, _, _]) Expr() (string, []ExprType, []error) {
	if len(sq.errs) != 0 {
		return "", nil, sq.errs
	}

	return sq.query, sq.args, nil
}

// TableExpr implements the TableExpr interface
func (sq *SubQueryExpr[TOuter, _, _]) TableExpr(TOuter) (string, []ExprType, []error) {
	return sq.Expr()
}

// TypedExpr implements the TypedExpr interface
func (sq *SubQueryExpr[_, _, S]) TypedExpr(S) (string, []ExprType, []error) {
	return sq.Expr()
}
