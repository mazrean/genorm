package genorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type SelectContext[T Table] struct {
	*Context[T]
	distinct        bool
	fields          []TableColumns[T]
	whereCondition  whereConditionClause[T]
	groupExpr       groupClause[T]
	havingCondition whereConditionClause[T]
	order           orderClause[T]
	limit           limitClause
	offset          offsetClause
	lockType        lockClause
}

func Select[T Table](table T) *SelectContext[T] {
	return &SelectContext[T]{
		Context: newContext(table),
	}
}

func (c *SelectContext[Table]) Distinct() *SelectContext[Table] {
	if c.distinct {
		c.addError(errors.New("distinct already set"))
		return c
	}

	c.distinct = true

	return c
}

func (c *SelectContext[Table]) Fields(fields ...TableColumns[Table]) *SelectContext[Table] {
	if c.fields != nil {
		c.addError(errors.New("fields already set"))
		return c
	}
	if len(fields) == 0 {
		c.addError(errors.New("no fields"))
		return c
	}

	fields = append(c.fields, fields...)
	fieldMap := make(map[TableColumns[Table]]struct{}, len(fields))
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

func (c *SelectContext[Table]) Where(
	condition TypedTableExpr[Table, WrappedPrimitive[bool]],
) *SelectContext[Table] {
	err := c.whereCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("where condition: %w", err))
	}

	return c
}

func (c *SelectContext[Table]) GroupBy(exprs ...TableExpr[Table]) *SelectContext[Table] {
	err := c.groupExpr.set(exprs)
	if err != nil {
		c.addError(fmt.Errorf("group by: %w", err))
	}

	return c
}

func (c *SelectContext[Table]) Having(
	condition TypedTableExpr[Table, WrappedPrimitive[bool]],
) *SelectContext[Table] {
	err := c.havingCondition.set(condition)
	if err != nil {
		c.addError(fmt.Errorf("having condition: %w", err))
	}

	return c
}

func (c *SelectContext[Table]) OrderBy(direction OrderDirection, expr TableExpr[Table]) *SelectContext[Table] {
	err := c.order.add(orderItem[Table]{
		expr:      expr,
		direction: direction,
	})
	if err != nil {
		c.addError(fmt.Errorf("order by: %w", err))
	}

	return c
}

func (c *SelectContext[Table]) Limit(limit uint64) *SelectContext[Table] {
	err := c.limit.set(limit)
	if err != nil {
		c.addError(fmt.Errorf("limit: %w", err))
	}

	return c
}

func (c *SelectContext[Table]) Offset(offset uint64) *SelectContext[Table] {
	err := c.offset.set(offset)
	if err != nil {
		c.addError(fmt.Errorf("offset: %w", err))
	}

	return c
}

func (c *SelectContext[Table]) Lock(lockType LockType) *SelectContext[Table] {
	err := c.lockType.set(lockType)
	if err != nil {
		c.addError(fmt.Errorf("lockType: %w", err))
	}

	return c
}

func (c *SelectContext[Table]) FindCtx(ctx context.Context, db DB) ([]Table, error) {
	errs := c.Errors()
	if len(errs) != 0 {
		return nil, errs[0]
	}

	columns, query, exprArgs, err := c.buildQuery()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	args := make([]any, 0, len(exprArgs))
	for _, arg := range exprArgs {
		args = append(args, arg)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return []Table{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	tables := []Table{}
	for rows.Next() {
		var table Table
		iTable := table.New()
		switch v := iTable.(type) {
		case Table:
			table = v
		default:
			return nil, fmt.Errorf("invalid table type: %T", iTable)
		}

		columnMap := table.ColumnMap()

		dests := make([]any, 0, len(columns))
		for _, column := range columns {
			columnField, ok := columnMap[column.SQLColumnName()]
			if !ok {
				return nil, fmt.Errorf("column %s not found", column.SQLColumnName())
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

func (c *SelectContext[Table]) Find(db DB) ([]Table, error) {
	return c.FindCtx(context.Background(), db)
}

func (c *SelectContext[Table]) TakeCtx(ctx context.Context, db DB) (Table, error) {
	var table Table

	err := c.limit.set(1)
	if err != nil {
		return table, fmt.Errorf("set limit 1: %w", err)
	}

	errs := c.Errors()
	if len(errs) != 0 {
		return table, errs[0]
	}

	columns, query, exprArgs, err := c.buildQuery()
	if err != nil {
		return table, fmt.Errorf("build query: %w", err)
	}

	args := make([]any, 0, len(exprArgs))
	for _, arg := range exprArgs {
		args = append(args, arg)
	}

	row := db.QueryRowContext(ctx, query, args...)

	iTable := table.New()
	switch v := iTable.(type) {
	case Table:
		table = v
	default:
		return table, fmt.Errorf("invalid table type: %T", iTable)
	}

	columnMap := table.ColumnMap()

	dests := make([]any, 0, len(columns))
	for _, column := range columns {
		columnField, ok := columnMap[column.SQLColumnName()]
		if !ok {
			return table, fmt.Errorf("column %s not found", column)
		}

		dests = append(dests, columnField)
	}

	err = row.Scan(dests...)
	if errors.Is(err, sql.ErrNoRows) {
		return table, ErrRecordNotFound
	}
	if err != nil {
		return table, fmt.Errorf("query: %w", err)
	}

	return table, nil
}

func (c *SelectContext[Table]) Take(db DB) (Table, error) {
	return c.TakeCtx(context.Background(), db)
}

func (c *SelectContext[Table]) buildQuery() ([]Column, string, []ExprType, error) {
	sb := strings.Builder{}
	args := []ExprType{}

	str := "SELECT "
	_, err := sb.WriteString(str)
	if err != nil {
		return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	if c.distinct {
		str = "DISTINCT "
		_, err = sb.WriteString(str)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}
	}

	var columns []Column
	if len(c.fields) == 0 {
		columns = c.table.Columns()
	} else {
		columns = make([]Column, 0, len(c.fields))
		for _, field := range c.fields {
			columns = append(columns, field)
		}
	}

	columnAliasMap := map[string]struct{}{}
	selectExprs := make([]string, 0, len(columns))
	for _, column := range columns {
		var alias string
		i := 0
		for ok := true; ok; _, ok = columnAliasMap[alias] {
			alias = fmt.Sprintf("%s_%s_%d", column.TableName(), column.ColumnName(), i)
			i++
		}

		columnAliasMap[alias] = struct{}{}
		selectExprs = append(selectExprs, fmt.Sprintf("%s AS %s", column.SQLColumnName(), alias))
	}

	str = strings.Join(selectExprs, ", ")
	_, err = sb.WriteString(str)
	if err != nil {
		return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	str = " FROM "
	_, err = sb.WriteString(str)
	if err != nil {
		return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	tableQuery, tableArgs, errs := c.table.Expr()
	if len(errs) != 0 {
		return nil, "", nil, fmt.Errorf("table expr: %w", errs[0])
	}

	_, err = sb.WriteString(tableQuery)
	if err != nil {
		return nil, "", nil, fmt.Errorf("write string(%s): %w", tableQuery, err)
	}

	args = append(args, tableArgs...)

	if c.whereCondition.exists() {
		whereQuery, whereArgs, err := c.whereCondition.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("where condition: %w", err)
		}

		str = " WHERE "
		_, err = sb.WriteString(str)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(whereQuery)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", whereQuery, err)
		}

		args = append(args, whereArgs...)
	}

	if c.groupExpr.exists() {
		groupExpr, groupArgs, err := c.groupExpr.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("group expr: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(groupExpr)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", groupExpr, err)
		}

		args = append(args, groupArgs...)
	}

	if c.havingCondition.exists() {
		havingQuery, havingArgs, err := c.havingCondition.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("having condition: %w", err)
		}

		str = " HAVING "
		_, err = sb.WriteString(str)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(havingQuery)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", havingQuery, err)
		}

		args = append(args, havingArgs...)
	}

	if c.order.exists() {
		orderQuery, orderArgs, err := c.order.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("order: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(orderQuery)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", orderQuery, err)
		}

		args = append(args, orderArgs...)
	}

	if c.limit.exists() {
		limitQuery, limitArgs, err := c.limit.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("limit: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(limitQuery)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", limitQuery, err)
		}

		args = append(args, limitArgs...)
	}

	if c.offset.exists() {
		offsetQuery, offsetArgs, err := c.offset.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("offset: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(offsetQuery)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", offsetQuery, err)
		}

		args = append(args, offsetArgs...)
	}

	if c.lockType.exists() {
		lockQuery, lockArgs, err := c.lockType.getExpr()
		if err != nil {
			return nil, "", nil, fmt.Errorf("lock: %w", err)
		}

		str = " "
		_, err = sb.WriteString(str)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", str, err)
		}

		_, err = sb.WriteString(lockQuery)
		if err != nil {
			return nil, "", nil, fmt.Errorf("write string(%s): %w", lockQuery, err)
		}

		args = append(args, lockArgs...)
	}

	return columns, sb.String(), args, nil
}
