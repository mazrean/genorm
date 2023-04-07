package relation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mazrean/genorm"
)

//nolint:revive
type RelationContext[S Table, T Table, _ JoinedTablePointer[V], V any] struct {
	baseTable S
	refTable  T
}

func NewRelationContext[S Table, T Table, U JoinedTablePointer[V], V any](baseTable S, refTable T) *RelationContext[S, T, U, V] {
	return &RelationContext[S, T, U, V]{
		baseTable: baseTable,
		refTable:  refTable,
	}
}

// Join INNER JOIN(CROSS JOIN)
func (r *RelationContext[S, T, U, V]) Join(
	expr genorm.TypedTableExpr[U, genorm.WrappedPrimitive[bool]],
) U {
	var joinedTable V

	relation, err := newRelation(join, r.baseTable, r.refTable, expr)
	if err != nil {
		U(&joinedTable).AddError(err)
		return &joinedTable
	}

	U(&joinedTable).SetRelation(relation)

	return &joinedTable
}

// LeftJoin LEFT JOIN
func (r *RelationContext[S, T, U, V]) LeftJoin(
	expr genorm.TypedTableExpr[U, genorm.WrappedPrimitive[bool]],
) U {
	var joinedTable V

	relation, err := newRelation(leftJoin, r.baseTable, r.refTable, expr)
	if err != nil {
		U(&joinedTable).AddError(err)
		return &joinedTable
	}

	U(&joinedTable).SetRelation(relation)

	return &joinedTable
}

// RightJoin RIGHT JOIN
func (r *RelationContext[S, T, U, V]) RightJoin(
	expr genorm.TypedTableExpr[U, genorm.WrappedPrimitive[bool]],
) U {
	var joinedTable V

	relation, err := newRelation(rightJoin, r.baseTable, r.refTable, expr)
	if err != nil {
		U(&joinedTable).AddError(err)
		return &joinedTable
	}

	U(&joinedTable).SetRelation(relation)

	return &joinedTable
}

type Relation struct {
	relationType RelationType
	baseTable    Table
	refTable     Table
	onExpr       genorm.Expr
}

func newRelation(relationType RelationType, baseTable, refTable Table, expr genorm.Expr) (*Relation, error) {
	if err := relationType.validate(); err != nil {
		return nil, fmt.Errorf("validate relation type: %w", err)
	}

	return &Relation{
		relationType: relationType,
		baseTable:    baseTable,
		refTable:     refTable,
		onExpr:       expr,
	}, nil
}

func (r *Relation) JoinedTableName() (string, []genorm.ExprType, []error) {
	sb := strings.Builder{}
	args := []genorm.ExprType{}

	str := "("
	_, err := sb.WriteString(str)
	if err != nil {
		return "", nil, []error{fmt.Errorf("write string(%s): %w", str, err)}
	}

	baseTableQuery, baseTableArgs, errs := r.baseTable.Expr()
	if len(errs) != 0 {
		return "", nil, errs
	}

	_, err = sb.WriteString(baseTableQuery)
	if err != nil {
		return "", nil, []error{fmt.Errorf("write string(%s): %w", baseTableQuery, err)}
	}

	args = append(args, baseTableArgs...)

	switch r.relationType {
	case join:
		if r.onExpr != nil {
			str = " INNER JOIN "
			_, err = sb.WriteString(str)
			if err != nil {
				return "", nil, []error{fmt.Errorf("write string(%s): %w", str, err)}
			}
		} else {
			str = " CROSS JOIN "
			_, err = sb.WriteString(str)
			if err != nil {
				return "", nil, []error{fmt.Errorf("write string(%s): %w", str, err)}
			}
		}
	case leftJoin:
		str = " LEFT JOIN "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, []error{fmt.Errorf("write string(%s): %w", str, err)}
		}
	case rightJoin:
		str = " RIGHT JOIN "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, []error{fmt.Errorf("write string(%s): %w", str, err)}
		}
	default:
		return "", nil, []error{errors.New("unsupported relation type")}
	}

	refTableQuery, refTableArgs, errs := r.refTable.Expr()
	if len(errs) != 0 {
		return "", nil, errs
	}

	_, err = sb.WriteString(refTableQuery)
	if err != nil {
		return "", nil, []error{fmt.Errorf("write string(%s): %w", refTableQuery, err)}
	}

	args = append(args, refTableArgs...)

	if r.onExpr != nil {
		str = " ON "
		_, err = sb.WriteString(str)
		if err != nil {
			return "", nil, []error{fmt.Errorf("write string(%s): %w", str, err)}
		}

		onExprQuery, onExprArgs, errs := r.onExpr.Expr()
		if len(errs) != 0 {
			return "", nil, errs
		}

		_, err = sb.WriteString(onExprQuery)
		if err != nil {
			return "", nil, []error{fmt.Errorf("write string(%s): %w", onExprQuery, err)}
		}

		args = append(args, onExprArgs...)
	}

	str = ")"
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, []error{fmt.Errorf("write string(%s): %w", str, err)}
	}

	return sb.String(), args, nil
}

//nolint:revive
type RelationType int8

const (
	join RelationType = iota + 1
	leftJoin
	rightJoin
)

func (rt RelationType) validate() error {
	if rt != join && rt != leftJoin && rt != rightJoin {
		return errors.New("unsupported relation type")
	}

	return nil
}
