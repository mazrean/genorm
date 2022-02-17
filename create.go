package genorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type CreateContext[T BasicTable] struct {
	*Context[T]
	values []T
	fields []TableColumns[T]
}

func NewCreateContext[T BasicTable](table T, tableBases ...T) *CreateContext[T] {
	ctx := newContext(table)
	if len(tableBases) == 0 {
		ctx.addError(errors.New("no insert values"))

		return &CreateContext[T]{
			Context: ctx,
		}
	}

	return &CreateContext[T]{
		Context: ctx,
		values:  tableBases,
		fields:  nil,
	}
}

func (c *CreateContext[BasicTable]) Fields(fields ...TableColumns[BasicTable]) *CreateContext[BasicTable] {
	if c.fields != nil {
		c.addError(errors.New("fields already set"))
		return c
	}
	if len(fields) == 0 {
		c.addError(errors.New("no fields"))
		return c
	}

	fields = append(c.fields, fields...)
	fieldMap := make(map[TableColumns[BasicTable]]struct{}, len(fields))
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

func (c *CreateContext[BasicTable]) DoContext(ctx context.Context, db *sql.DB) (rowsAffected int64, err error) {
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

func (c *CreateContext[BasicTable]) Do(db *sql.DB) (rowsAffected int64, err error) {
	return c.DoContext(context.Background(), db)
}

func (c *CreateContext[BasicTable]) buildQuery() (string, []any, error) {
	args := []any{}

	sb := strings.Builder{}
	sb.WriteString("INSERT INTO ")
	sb.WriteString(c.table.TableName())

	var fields []string
	if c.fields == nil {
		columns := c.table.Columns()
		fields = make([]string, 0, len(columns))
		for _, column := range columns {
			fields = append(fields, column.ColumnName())
		}
	} else {
		fields = make([]string, 0, len(c.fields))
		for _, field := range c.fields {
			fields = append(fields, field.ColumnName())
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
		basicTable, ok := val.(BasicTable)
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

func (c *CreateContext[TableBase]) buildValueList(sb strings.Builder, args []any, fields []string, fieldValueMap map[string]ColumnField) (strings.Builder, []any, error) {
	sb.WriteString("(")
	for i, columnName := range fields {
		if i != 0 {
			sb.WriteString(", ")
		}

		columnField, ok := fieldValueMap[columnName]
		if !ok {
			return sb, nil, errors.New("field not found")
		}

		fieldValue, err := columnField.iValue()
		if err != nil && !errors.Is(err, ErrNullValue) {
			return sb, nil, fmt.Errorf("failed to get field value: %w", err)
		}

		if errors.Is(err, ErrNullValue) {
			sb.WriteString("NULL")
		} else {
			sb.WriteString("?")
			args = append(args, fieldValue)
		}
	}
	sb.WriteString(")")

	return sb, args, nil
}
