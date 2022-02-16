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
		values:  tableBases,
		fields:  nil,
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
	var v interface{} = c.table
	if basicTable, ok := v.(BasicTable); ok {
		return c.basicTableBuildQuery(basicTable)
	} else if joinedTable, ok := v.(JoinedTable); ok {
		return c.joinedTableBuildQuery(joinedTable)
	}

	return "", nil, errors.New("invalid table")
}

func (c *CreateContext[TableBase]) basicTableBuildQuery(basicTable BasicTable) (string, []any, error) {
	args := []any{}

	sb := strings.Builder{}
	sb.WriteString("INSERT INTO ")
	sb.WriteString(c.table.TableName())

	var fields []string
	sb.WriteString(" (")
	if c.fields == nil {
		fields = c.basicTableFields(basicTable)
		sb.WriteString(strings.Join(fields, ", "))
	} else {
		for i, field := range c.fields {
			if i != 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(field.ColumnName())

			fields = append(fields, field.ColumnName())
		}
	}
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

		fieldValueMap, err := c.basicTableFieldValueMap(basicTable)
		if err != nil {
			return "", nil, fmt.Errorf("basic table field map: %w", err)
		}

		sb, args, err = c.buildValueList(sb, args, fields, fieldValueMap)
		if err != nil {
			return "", nil, fmt.Errorf("build value list: %w", err)
		}
	}

	return sb.String(), args, nil
}

func (c *CreateContext[TableBase]) joinedTableBuildQuery(joinedTable JoinedTable) (string, []any, error) {
	args := []any{}

	basicTables := joinedTable.BaseTables()

	sb := strings.Builder{}
	sb.WriteString("INSERT INTO ")
	sb.WriteString(c.table.TableName())

	var fields []string
	sb.WriteString(" (")
	if c.fields == nil {
		fields = []string{}
		for _, basicTable := range basicTables {
			fields = append(fields, c.basicTableFields(basicTable)...)
		}
		sb.WriteString(strings.Join(fields, ", "))
	} else {
		for i, field := range c.fields {
			if i != 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(field.ColumnName())

			fields = append(fields, field.ColumnName())
		}
	}
	sb.WriteString(") VALUES ")

	for i, value := range c.values {
		if i != 0 {
			sb.WriteString(", ")
		}

		var val any = value
		joinedTable, ok := val.(JoinedTable)
		if !ok {
			return "", nil, errors.New("failed to cast value to joined table")
		}

		basicTables := joinedTable.BaseTables()

		fieldValueMap := make(map[string]any)
		for _, basicTable := range basicTables {
			basicTableFieldValueMap, err := c.basicTableFieldValueMap(basicTable)
			if err != nil {
				return "", nil, fmt.Errorf("basic table field map: %w", err)
			}

			for k, v := range basicTableFieldValueMap {
				if _, ok := fieldValueMap[k]; ok {
					return "", nil, fmt.Errorf("duplicate field: %s", k)
				}

				fieldValueMap[k] = v
			}
		}

		var err error
		sb, args, err = c.buildValueList(sb, args, fields, fieldValueMap)
		if err != nil {
			return "", nil, fmt.Errorf("build value list: %w", err)
		}
	}

	return sb.String(), args, nil
}

func (c *CreateContext[TableBase]) buildValueList(sb strings.Builder, args []any, fields []string, fieldValueMap map[string]any) (strings.Builder, []any, error) {
	sb.WriteString("(")
	for i, columnName := range fields {
		if i != 0 {
			sb.WriteString(", ")
		}

		iField, ok := fieldValueMap[columnName]
		if !ok {
			return sb, nil, errors.New("field not found")
		}

		columnField, ok := iField.(ColumnField)
		if !ok {
			return sb, nil, errors.New("invalid field type")
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

func (c *CreateContext[TableBase]) basicTableFields(basicTable BasicTable) []string {
	fields := make([]string, 0, len(basicTable.ColumnNames()))
	for _, field := range basicTable.ColumnNames() {
		fields = append(fields, fmt.Sprintf("%s.%s", basicTable.TableName(), field))
	}

	return fields
}

func (c *CreateContext[TableBase]) basicTableFieldValueMap(basicTable BasicTable) (map[string]any, error) {
	t := reflect.TypeOf(basicTable)
	v := reflect.ValueOf(basicTable)

	if t.Kind() != reflect.Struct {
		return nil, errors.New("value is not a struct")
	}

	fieldMap := make(map[string]any)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columnName, ok := field.Tag.Lookup("genorm")
		if !ok {
			return nil, errors.New("genorm tag not found")
		}

		fieldMap[fmt.Sprintf("%s.%s", basicTable.TableName(), columnName)] = v.Field(i).Interface()
	}

	return fieldMap, nil
}
