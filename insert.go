package genorm

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type InsertContext[T BasicTable] struct {
	*Context[T]
	values []T
	fields []TableColumns[T]
}

func Insert[T BasicTable](table T) *InsertContext[T] {
	ctx := newContext(table)

	return &InsertContext[T]{
		Context: ctx,
		fields:  nil,
	}
}

func (c *InsertContext[Table]) Values(tableBases ...Table) *InsertContext[Table] {
	if len(tableBases) == 0 {
		c.addError(errors.New("no values"))

		return c
	}
	if len(c.values) != 0 {
		c.addError(errors.New("values already set"))

		return c
	}

	c.values = append(c.values, tableBases...)

	return c
}

func (c *InsertContext[Table]) Fields(fields ...TableColumns[Table]) *InsertContext[Table] {
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

func (c *InsertContext[Table]) DoCtx(ctx context.Context, db DB) (rowsAffected int64, err error) {
	errs := c.Errors()
	if len(errs) != 0 {
		return 0, errs[0]
	}

	query, args, err := c.buildQuery()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (c *InsertContext[Table]) Do(db DB) (rowsAffected int64, err error) {
	return c.DoCtx(context.Background(), db)
}

func (c *InsertContext[Table]) buildQuery() (string, []any, error) {
	args := []any{}

	sb := &strings.Builder{}
	sb.WriteString("INSERT INTO ")
	sb.WriteString(c.table.TableName())

	var fields []string
	if c.fields == nil {
		columns := c.table.Columns()
		fields = make([]string, 0, len(columns))
		for _, column := range columns {
			fields = append(fields, column.SQLColumnName())
		}
	} else {
		fields = make([]string, 0, len(c.fields))
		for _, field := range c.fields {
			fields = append(fields, field.SQLColumnName())
		}
	}

	sb.WriteString(" (")
	sb.WriteString(strings.Join(fields, ", "))
	sb.WriteString(") VALUES ")

	for i, value := range c.values {
		if i != 0 {
			sb.WriteString(", ")
		}

		var val any = value
		basicTable, ok := val.(Table)
		if !ok {
			return "", nil, errors.New("failed to cast value to basic table")
		}

		var err error
		sb, args, err = c.buildValueList(sb, args, fields, basicTable.ColumnMap())
		if err != nil {
			return "", nil, fmt.Errorf("build value list: %w", err)
		}
	}

	return sb.String(), args, nil
}

func (c *InsertContext[Table]) buildValueList(sb *strings.Builder, args []any, fields []string, fieldValueMap map[string]ColumnFieldExprType) (*strings.Builder, []any, error) {
	sb.WriteString("(")
	for i, columnName := range fields {
		if i != 0 {
			sb.WriteString(", ")
		}

		columnField, ok := fieldValueMap[columnName]
		if !ok {
			return sb, nil, fmt.Errorf("field(%s) not found", columnName)
		}

		fieldValue, err := columnField.Value()
		if err != nil {
			return sb, nil, fmt.Errorf("failed to get field value: %w", err)
		}

		if fieldValue == nil {
			sb.WriteString("NULL")
		} else {
			sb.WriteString("?")
			args = append(args, fieldValue)
		}
	}
	sb.WriteString(")")

	return sb, args, nil
}
