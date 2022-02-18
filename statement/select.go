package statement

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mazrean/genorm"
)

type SelectContext[T Table] struct {
	*Context[T]
	distinct        bool
	fields          []genorm.TableColumns[T]
	whereCondition  whereConditionClause[T]
	groupExpr       []genorm.TableExpr[T]
	havingCondition whereConditionClause[T]
	order           orderClause[T]
	limit           limitClause
	offset          offsetClause
	lockType        LockType
}

func NewSelectContext[T Table](table T) *SelectContext[T] {
	return &SelectContext[T]{
		Context: newContext(table),
	}
}

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

func (c *SelectContext[TableBase]) Fields(fields ...genorm.TableColumns[TableBase]) *SelectContext[TableBase] {
	if c.fields != nil {
		c.addError(errors.New("fields already set"))
		return c
	}
	if len(fields) == 0 {
		c.addError(errors.New("no fields"))
		return c
	}

	fields = append(c.fields, fields...)
	fieldMap := make(map[genorm.TableColumns[TableBase]]struct{}, len(fields))
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

func (c *SelectContext[TableBase]) Where(condition genorm.TypedTableExpr[TableBase, bool]) *SelectContext[TableBase] {
	err := c.whereCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("where condition: %w", err))
	}

	return c
}

func (c *SelectContext[TableBase]) GroupBy(exprs ...genorm.TableExpr[TableBase]) *SelectContext[TableBase] {
	if len(exprs) == 0 {
		c.addError(errors.New("no group expr"))
		return c
	}

	c.groupExpr = append(c.groupExpr, exprs...)

	return c
}

func (c *SelectContext[TableBase]) Having(condition genorm.TypedTableExpr[TableBase, bool]) *SelectContext[TableBase] {
	err := c.havingCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("having condition: %w", err))
	}

	return c
}

func (c *SelectContext[TableBase]) OrderBy(direction OrderDirection, expr genorm.TableExpr[TableBase]) *SelectContext[TableBase] {
	err := c.order.add(orderItem[TableBase]{
		expr:      expr,
		direction: direction,
	})
	if err != nil {
		c.addError(fmt.Errorf("order by: %w", err))
	}

	return c
}

func (c *SelectContext[TableBase]) Limit(limit uint64) *SelectContext[TableBase] {
	err := c.limit.set(limit)
	if err != nil {
		c.addError(fmt.Errorf("limit: %w", err))
	}

	return c
}

func (c *SelectContext[TableBase]) Offset(offset uint64) *SelectContext[TableBase] {
	err := c.offset.set(offset)
	if err != nil {
		c.addError(fmt.Errorf("offset: %w", err))
	}

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
		columnMap := table.ColumnMap()

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

	var columns []genorm.Column
	if len(c.fields) == 0 {
		columns = c.table.Columns()
	} else {
		columns = make([]genorm.Column, 0, len(c.fields))
		for _, field := range c.fields {
			columns = append(columns, field)
		}
	}

	selectExprs := make([]string, 0, len(columns))
	for _, column := range columns {
		var alias string
		i := 0
		for ok := false; ok; _, ok = columnAliasMap[alias] {
			alias = fmt.Sprintf("%s_%s_%d", column.TableName(), column.ColumnName(), i)
			i++
		}

		columnAliasMap[alias] = column.SQLColumnName()
		selectExprs = append(selectExprs, fmt.Sprintf("%s AS %s", column.SQLColumnName(), alias))
	}

	sb.WriteString(strings.Join(selectExprs, ", "))

	sb.WriteString(" FROM ")
	sb.WriteString(c.table.SQLTableName())

	if c.whereCondition.exists() {
		whereQuery, whereArgs, err := c.whereCondition.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("where condition: %w", err)
		}

		sb.WriteString(" WHERE ")
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

	if c.havingCondition.exists() {
		havingQuery, havingArgs, err := c.havingCondition.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("having condition: %w", err)
		}

		sb.WriteString(" HAVING ")
		sb.WriteString(havingQuery)
		args = append(args, havingArgs...)
	}

	if c.order.exists() {
		orderQuery, orderArgs, err := c.order.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("order: %w", err)
		}

		sb.WriteString(" ")
		sb.WriteString(orderQuery)
		args = append(args, orderArgs...)
	}

	if c.limit.exists() {
		limitQuery, limitArgs, err := c.limit.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("limit: %w", err)
		}

		sb.WriteString(" ")
		sb.WriteString(limitQuery)
		args = append(args, limitArgs...)
	}

	if c.offset.exists() {
		offsetQuery, offsetArgs, err := c.offset.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("offset: %w", err)
		}

		sb.WriteString(" ")
		sb.WriteString(offsetQuery)
		args = append(args, offsetArgs...)
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
