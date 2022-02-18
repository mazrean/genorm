package genorm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mazrean/genorm"
)

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

func (r *Relation) JoinedTableName() (string, []any, error) {
	sb := strings.Builder{}
	args := []any{}

	sb.WriteString("(")

	baseTableQuery, baseTableArgs := r.baseTable.Expr()
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

	refTableQuery, refTableArgs := r.refTable.Expr()
	sb.WriteString(refTableQuery)
	args = append(args, refTableArgs...)

	if r.onExpr != nil {
		sb.WriteString(" ON ")

		onExprQuery, onExprArgs := r.onExpr.Expr()
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
