package relation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mazrean/genorm"
)

type RelationContext[BaseTable Table, RefTable Table, _ JoinedTable] struct {
	baseTable BaseTable
	refTable  RefTable
}

func NewRelationContext[S Table, T Table, U JoinedTable](baseTable S, refTable T) *RelationContext[S, T, U] {
	return &RelationContext[S, T, U]{
		baseTable: baseTable,
		refTable:  refTable,
	}
}

func (r *RelationContext[BaseTable, RefTable, JoinedTable]) Join(
	expr genorm.TypedTableExpr[JoinedTable, *genorm.WrappedPrimitive[bool]],
) JoinedTable {
	var (
		baseTable   BaseTable
		refTable    RefTable
		joinedTable JoinedTable
	)

	relation, err := newRelation(join, baseTable, refTable, expr)
	if err != nil {
		joinedTable.AddError(err)
		return joinedTable
	}

	joinedTable.SetRelation(relation)

	return joinedTable
}

func (r *RelationContext[BaseTable, RefTable, JoinedTable]) LeftJoin(
	expr genorm.TypedTableExpr[JoinedTable, *genorm.WrappedPrimitive[bool]],
) JoinedTable {
	var (
		baseTable   BaseTable
		refTable    RefTable
		joinedTable JoinedTable
	)

	relation, err := newRelation(leftJoin, baseTable, refTable, expr)
	if err != nil {
		joinedTable.AddError(err)
		return joinedTable
	}

	joinedTable.SetRelation(relation)

	return joinedTable
}

func (r *RelationContext[BaseTable, RefTable, JoinedTable]) RightJoin(
	expr genorm.TypedTableExpr[JoinedTable, *genorm.WrappedPrimitive[bool]],
) JoinedTable {
	var (
		baseTable   BaseTable
		refTable    RefTable
		joinedTable JoinedTable
	)

	relation, err := newRelation(rightJoin, baseTable, refTable, expr)
	if err != nil {
		joinedTable.AddError(err)
		return joinedTable
	}

	joinedTable.SetRelation(relation)

	return joinedTable
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

func (r *Relation) JoinedTableName() (string, []genorm.ExprType, error) {
	sb := strings.Builder{}
	args := []genorm.ExprType{}

	sb.WriteString("(")

	baseTableQuery, baseTableArgs, errs := r.baseTable.Expr()
	if len(errs) != 0 {
		return "", nil, errs[0]
	}

	sb.WriteString(baseTableQuery)
	args = append(args, baseTableArgs...)

	switch r.relationType {
	case join:
		if r.onExpr != nil {
			sb.WriteString(" INNER JOIN ")
		} else {
			sb.WriteString(" CROSS JOIN ")
		}
	case leftJoin:
		sb.WriteString(" LEFT JOIN ")
	case rightJoin:
		sb.WriteString(" RIGHT JOIN ")
	default:
		return "", nil, errors.New("unsupported relation type")
	}

	refTableQuery, refTableArgs, errs := r.refTable.Expr()
	if len(errs) != 0 {
		return "", nil, errs[0]
	}

	sb.WriteString(refTableQuery)
	args = append(args, refTableArgs...)

	if r.onExpr != nil {
		sb.WriteString(" ON ")

		onExprQuery, onExprArgs, errs := r.onExpr.Expr()
		if len(errs) != 0 {
			return "", nil, errs[0]
		}

		sb.WriteString(onExprQuery)
		args = append(args, onExprArgs...)
	}

	sb.WriteString(")")

	return sb.String(), args, nil
}

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
