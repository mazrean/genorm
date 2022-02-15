package genorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type CreateContext[T TableBase] struct {
	*Context[T]
	values []T
	fields []TableColumns[T]
}

func newCreateContext[T TableBase](ctx *Context[T], tableBases ...T) *CreateContext[T] {
	return &CreateContext[T]{
		Context: ctx,
		values: tableBases,
		fields: nil,
	}
}

func (c *CreateContext[TableBase]) Fields(fields ...TableColumns[TableBase]) *CreateContext[TableBase] {
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

func (c *CreateContext[TableBase]) DoContext(ctx context.Context, db *sql.DB) (rowsAffected int64, err error) {
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

func (c *CreateContext[TableBase]) Do(db *sql.DB) (rowsAffected int64, err error) {
	errs := c.Errors()
	if len(errs) != 0 {
		return 0, errs[0]
	}

	query, args, err := c.buildQuery()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (c *CreateContext[TableBase]) buildQuery() (string, []any, error) {
	args := []any{}

	sb := strings.Builder{}
	sb.WriteString("INSERT INTO ")
	sb.WriteString(c.table.TableName())

	var v interface{} = c.table
	basicTable, ok := v.(BasicTable)
	if ok {
		var fields []string
		sb.WriteString(" (")
		if c.fields == nil {
			fields = make([]string, 0, len(basicTable.ColumnNames()))
			for i, field := range basicTable.ColumnNames() {
				if i != 0 {
					sb.WriteString(", ")
				}

				sb.WriteString(basicTable.TableName())
				sb.WriteString(".")
				sb.WriteString(field)

				fields = append(fields, fmt.Sprintf("%s.%s", basicTable.TableName(), field))
			}
		} else {
			for i, field := range c.fields {
				if i != 0 {
					sb.WriteString(", ")
				}

				sb.WriteString(string(field))

				fields = append(fields, string(field))
			}
		}
		sb.WriteString(") VALUES ")

		for i, value := range c.values {
			if i != 0 {
				sb.WriteString(", ")
			}

			t := reflect.TypeOf(value)
			v := reflect.ValueOf(value)

			if t.Kind() != reflect.Struct {
				return "", nil, errors.New("value is not a struct")
			}

			fieldMap := make(map[string]reflect.StructField)
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				columnName := field.Tag.Get("genorm")

				fieldMap[fmt.Sprintf("%s.%s", basicTable.TableName(), columnName)] = field
			}

			sb.WriteString("(")
			for i, columnName := range fields {
				if i != 0 {
					sb.WriteString(", ")
				}

				field, ok := fieldMap[columnName]
				if !ok {
					return "", nil, errors.New("field not found")
				}

				iField := v.Elem().FieldByName(field.Name).Interface()
				columnField, ok := iField.(ColumnField)
				if !ok {
					return "", nil, errors.New("invalid field type")
				}

				fieldValue, err := columnField.iValue()
				if err != nil && !errors.Is(err, ErrNullValue) {
					return "", nil, fmt.Errorf("failed to get field value: %w", err)
				}

				if errors.Is(err, ErrNullValue) {
					sb.WriteString("NULL")
				} else {
					sb.WriteString("?")
					args = append(args, fieldValue)
				}
			}
			sb.WriteString(")")
		}
	}

	return sb.String(), args, nil
}
