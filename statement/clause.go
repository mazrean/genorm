package statement

import (
	"errors"

	"github.com/mazrean/genorm"
)

type whereConditionClause[T Table] struct {
	condition genorm.TypedTableExpr[T, bool]
}

func (c *whereConditionClause[Table]) set(condition genorm.TypedTableExpr[Table, bool]) error {
	if c.condition != nil {
		return errors.New("where conditions already set")
	}
	if condition == nil {
		return errors.New("empty where condition")
	}

	c.condition = condition

	return nil
}

func (c *whereConditionClause[Table]) exists() bool {
	return c.condition != nil
}

func (c *whereConditionClause[Table]) getExpr() (string, []any, error) {
	query, args := c.condition.Expr()

	return query, args, nil
}
