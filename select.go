package genorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type SelectContext[T TableBase] struct {
	*Context[T]
	distinct        bool
	fields          []TableColumns[T]
	whereCondition  TypedTableExpr[T, bool]
	groupExpr       []TableExpr[T]
	havingCondition TypedTableExpr[T, bool]
	orderExprs      []orderItem[T]
	limit           uint64
	offset          uint64
	lockType        LockType
}

func NewSelectContext[T TableBase](table T) *SelectContext[T] {
	return &SelectContext[T]{
		Context: newContext(table),
	}
}

type orderItem[T TableBase] struct {
	expr      TableExpr[T]
	direction OrderDirection
}

type OrderDirection uint8

const (
	Asc OrderDirection = iota + 1
	Desc
)

type LockType uint8

const (
	none LockType = iota
	ForUpdate
	ForShare
)

func (c *SelectContext[TableBase]) Distinct() *SelectContext[TableBase] {
	if c.distinct {
		c.addError(errors.New("distinct already set"))
		return c
	}

	c.distinct = true

	return c
}

func (c *SelectContext[TableBase]) Fields(fields ...TableColumns[TableBase]) *SelectContext[TableBase] {
	if c.fields != nil {
		c.addError(errors.New("fields already set"))
		return c
	}
	if len(fields) == 0 {
		c.addError(errors.New("no fields"))
		return c
	}

	fields = append(c.fields, fields...)
	fieldMap := make(map[TableColumns[TableBase]]struct{}, len(fields))
	for _, field := range fields {
		if _, ok := fieldMap[field]; ok {
			c.addError(errors.New("duplicate field"))
			return c
		}

		fieldMap[field] = struct{}{}
	}

	c.fields = fields

	return c
}

func (c *SelectContext[TableBase]) Where(condition TypedTableExpr[TableBase, bool]) *SelectContext[TableBase] {
	if c.whereCondition != nil {
		c.addError(errors.New("where conditions already set"))
		return c
	}
	if condition == nil {
		c.addError(errors.New("empty where condition"))
		return c
	}

	c.whereCondition = condition

	return c
}

func (c *SelectContext[TableBase]) GroupBy(exprs ...TableExpr[TableBase]) *SelectContext[TableBase] {
	if len(exprs) == 0 {
		c.addError(errors.New("no group expr"))
		return c
	}

	c.groupExpr = append(c.groupExpr, exprs...)

	return c
}

func (c *SelectContext[TableBase]) Having(condition TypedTableExpr[TableBase, bool]) *SelectContext[TableBase] {
	if c.havingCondition != nil {
		c.addError(errors.New("having conditions already set"))
		return c
	}
	if condition == nil {
		c.addError(errors.New("empty having condition"))
		return c
	}

	c.havingCondition = condition

	return c
}

func (c *SelectContext[TableBase]) OrderBy(direction OrderDirection, expr TableExpr[TableBase]) *SelectContext[TableBase] {
	if expr == nil {
		c.addError(errors.New("empty order expr"))
		return c
	}

	if direction != Asc && direction != Desc {
		c.addError(errors.New("invalid order direction"))
		return c
	}

	c.orderExprs = append(c.orderExprs, orderItem[TableBase]{
		expr:      expr,
		direction: direction,
	})

	return c
}

func (c *SelectContext[TableBase]) Limit(limit uint64) *SelectContext[TableBase] {
	if c.limit != 0 {
		c.addError(errors.New("limit already set"))
		return c
	}
	if limit == 0 {
		c.addError(errors.New("invalid limit"))
		return c
	}

	c.limit = limit

	return c
}

func (c *SelectContext[TableBase]) Offset(offset uint64) *SelectContext[TableBase] {
	if c.offset != 0 {
		c.addError(errors.New("offset already set"))
		return c
	}

	c.offset = offset

	return c
}

func (c *SelectContext[TableBase]) Lock(lockType LockType) *SelectContext[TableBase] {
	if c.lockType != none {
		c.addError(errors.New("lock already set"))
		return c
	}

	if lockType != ForUpdate && lockType != ForShare {
		c.addError(errors.New("invalid lock type"))
		return c
	}

	c.lockType = lockType

	return c
}

func (c *SelectContext[TableBase]) DoContext(ctx context.Context, db *sql.DB) ([]TableBase, error) {
	errs := c.Errors()
	if len(errs) != 0 {
		return nil, errs[0]
	}

	columnAliasMap, query, args, err := c.buildQuery()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	resultColumns, err := rows.Columns()

	tables := []TableBase{}
	for rows.Next() {
		var table TableBase

		var v any = table
		var columnMap map[string]ColumnField
		if basicTable, ok := v.(BasicTable); ok {
			columnMap = basicTable.ColumnMap()
		} else if joinedTable, ok := v.(JoinedTable); ok {
			columnMap = map[string]ColumnField{}
			for _, basicTable := range joinedTable.BaseTables() {
				basicColumnMap := basicTable.ColumnMap()
				for columnID, field := range basicColumnMap {
					columnMap[columnID] = field
				}
			}
		} else {
			return nil, fmt.Errorf("invalid table type: %T", v)
		}

		dests := make([]any, 0, len(resultColumns))
		for _, columnAlias := range resultColumns {
			columnField, ok := columnMap[columnAliasMap[columnAlias]]
			if !ok {
				return nil, fmt.Errorf("column %s not found", columnAlias)
			}

			dests = append(dests, columnField)
		}

		err := rows.Scan(dests...)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func (c *SelectContext[TableBase]) Do(db *sql.DB) ([]TableBase, error) {
	return c.DoContext(context.Background(), db)
}

func (c *SelectContext[TableBase]) buildQuery() (map[string]string, string, []any, error) {
	columnAliasMap := map[string]string{}
	sb := strings.Builder{}
	args := []any{}

	sb.WriteString("SELECT ")

	if c.distinct {
		sb.WriteString("DISTINCT ")
	}

	var fields []string
	if len(c.fields) == 0 {
		c.table.Columns()
	} else {
		fields = make([]string, 0, len(c.fields))
		for _, field := range c.fields {
			fields = append(fields, field.ColumnID())
		}
	}

	sb.WriteString(strings.Join(fields, ", "))

	sb.WriteString(" FROM ")
	sb.WriteString(c.table.TableID())

	if c.whereCondition != nil {
		sb.WriteString(" WHERE ")

		whereQuery, whereArgs := c.whereCondition.Expr()
		sb.WriteString(whereQuery)
		args = append(args, whereArgs...)
	}

	if len(c.groupExpr) != 0 {
		sb.WriteString(" GROUP BY ")

		groupQueries := make([]string, 0, len(c.groupExpr))
		for _, groupExpr := range c.groupExpr {
			groupQuery, groupArgs := groupExpr.Expr()

			groupQueries = append(groupQueries, groupQuery)
			args = append(args, groupArgs...)
		}

		sb.WriteString(strings.Join(groupQueries, ", "))
	}

	if c.havingCondition != nil {
		sb.WriteString(" HAVING ")

		havingQuery, havingArgs := c.havingCondition.Expr()
		sb.WriteString(havingQuery)
		args = append(args, havingArgs...)
	}

	if len(c.orderExprs) != 0 {
		sb.WriteString(" ORDER BY ")

		orderQueries := make([]string, 0, len(c.orderExprs))
		for _, orderItem := range c.orderExprs {
			orderQuery, orderArgs := orderItem.expr.Expr()

			var directionQuery string
			switch orderItem.direction {
			case Asc:
				directionQuery = "ASC"
			case Desc:
				directionQuery = "DESC"
			default:
				return nil, "", nil, fmt.Errorf("invalid order direction: %d", orderItem.direction)
			}

			orderQueries = append(orderQueries, fmt.Sprintf("%s %s", orderQuery, directionQuery))
			args = append(args, orderArgs...)
		}

		sb.WriteString(strings.Join(orderQueries, ", "))
	}

	if c.limit != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", c.limit))
	}

	if c.offset != 0 {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", c.offset))
	}

	if c.lockType != none {
		switch c.lockType {
		case ForUpdate:
			sb.WriteString(" FOR UPDATE")
		case ForShare:
			sb.WriteString(" FOR SHARE")
		}
	}

	return columnAliasMap, sb.String(), args, nil
}
