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

func (c *InsertContext[T]) Values(tableBases ...T) *InsertContext[T] {
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

func (c *InsertContext[T]) Fields(fields ...TableColumns[T]) *InsertContext[T] {
	if c.fields != nil {
		c.addError(errors.New("fields already set"))
		return c
	}
	if len(fields) == 0 {
		c.addError(errors.New("no fields"))
		return c
	}

	fields = append(c.fields, fields...)
	fieldMap := make(map[TableColumns[T]]struct{}, len(fields))
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

func (c *InsertContext[T]) DoCtx(ctx context.Context, db DB) (rowsAffected int64, err error) {
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

func (c *InsertContext[T]) Do(db DB) (rowsAffected int64, err error) {
	return c.DoCtx(context.Background(), db)
}

func (c *InsertContext[T]) buildQuery() (string, []any, error) {
	args := []any{}

	sb := &strings.Builder{}

	str := "INSERT INTO `"
	_, err := sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	str = c.table.TableName()
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	str = "`"
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

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

	str = " ("
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	str = strings.Join(fields, ", ")
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	str = ") VALUES "
	_, err = sb.WriteString(str)
	if err != nil {
		return "", nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	for i, value := range c.values {
		if i != 0 {
			str = ", "
			_, err = sb.WriteString(str)
			if err != nil {
				return "", nil, fmt.Errorf("write string(%s): %w", str, err)
			}
		}

		var err error
		sb, args, err = c.buildValueList(sb, args, fields, value.ColumnMap())
		if err != nil {
			return "", nil, fmt.Errorf("build value list: %w", err)
		}
	}

	return sb.String(), args, nil
}

func (c *InsertContext[T]) buildValueList(sb *strings.Builder, args []any, fields []string, fieldValueMap map[string]ColumnFieldExprType) (*strings.Builder, []any, error) {
	str := "("
	_, err := sb.WriteString(str)
	if err != nil {
		return sb, args, fmt.Errorf("write string(%s): %w", str, err)
	}

	for i, columnName := range fields {
		if i != 0 {
			str = ", "
			_, err = sb.WriteString(str)
			if err != nil {
				return sb, args, fmt.Errorf("write string(%s): %w", str, err)
			}
		}

		columnField, ok := fieldValueMap[columnName]
		if !ok {
			return sb, nil, fmt.Errorf("field(%s) not found", columnName)
		}

		_, err := columnField.Value()
		if err != nil && !errors.Is(err, ErrNullValue) {
			return sb, nil, fmt.Errorf("failed to get field value: %w", err)
		}

		if errors.Is(err, ErrNullValue) {
			str = "NULL"
			_, err = sb.WriteString(str)
			if err != nil {
				return sb, nil, fmt.Errorf("write string(%s): %w", str, err)
			}
		} else {
			str = "?"
			_, err = sb.WriteString(str)
			if err != nil {
				return sb, nil, fmt.Errorf("write string(%s): %w", str, err)
			}

			args = append(args, columnField)
		}
	}

	str = ")"
	_, err = sb.WriteString(str)
	if err != nil {
		return sb, nil, fmt.Errorf("write string(%s): %w", str, err)
	}

	return sb, args, nil
}
