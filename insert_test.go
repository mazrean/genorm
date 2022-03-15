package genorm_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestInsertBuildQuery(t *testing.T) {
	t.Parallel()

	type field struct {
		tableName     string
		columnName    string
		sqlColumnName string
	}

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	type orderItem struct {
		direction genorm.OrderDirection
		expr      expr
	}

	columnFieldExpr1 := genorm.Wrap(1)
	columnFieldExpr2 := genorm.Wrap(2)

	tests := []struct {
		description string
		tableName   string
		isFieldSet  bool
		fields      []string
		values      []map[string]genorm.ColumnFieldExprType
		query       string
		args        []any
		err         bool
	}{
		{
			description: "normal",
			tableName:   "hoge",
			fields:      []string{"`hoge`.`huga`"},
			values: []map[string]genorm.ColumnFieldExprType{
				{
					"`hoge`.`huga`": &columnFieldExpr1,
				},
			},
			query: "INSERT INTO `hoge` (`hoge`.`huga`) VALUES (?)",
			args:  []any{&columnFieldExpr1},
		},
		{
			description: "multi fields",
			tableName:   "hoge",
			fields:      []string{"`hoge`.`huga`", "`hoge`.`piyo`"},
			values: []map[string]genorm.ColumnFieldExprType{
				{
					"`hoge`.`huga`": &columnFieldExpr1,
					"`hoge`.`piyo`": &columnFieldExpr2,
				},
			},
			query: "INSERT INTO `hoge` (`hoge`.`huga`, `hoge`.`piyo`) VALUES (?, ?)",
			args:  []any{&columnFieldExpr1, &columnFieldExpr2},
		},
		{
			description: "multi values",
			tableName:   "hoge",
			fields:      []string{"`hoge`.`huga`"},
			values: []map[string]genorm.ColumnFieldExprType{
				{
					"`hoge`.`huga`": &columnFieldExpr1,
				},
				{
					"`hoge`.`huga`": &columnFieldExpr2,
				},
			},
			query: "INSERT INTO `hoge` (`hoge`.`huga`) VALUES (?), (?)",
			args:  []any{&columnFieldExpr1, &columnFieldExpr2},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			table := mock.NewMockBasicTable(ctrl)
			table.
				EXPECT().
				TableName().
				Return(test.tableName)
			table.
				EXPECT().
				GetErrors().
				Return(nil)

			builder := genorm.Insert(table)

			fields := make([]genorm.Column, 0, len(test.fields))
			if test.isFieldSet {
				tableFields := make([]genorm.TableColumns[*mock.MockBasicTable], 0, len(test.fields))
				for _, field := range test.fields {
					mockColumn := mock.NewMockTypedTableColumn[*mock.MockBasicTable, genorm.WrappedPrimitive[bool]](ctrl)
					mockColumn.
						EXPECT().
						SQLColumnName().
						Return(field)

					tableFields = append(tableFields, mockColumn)
					fields = append(fields, mockColumn)
				}

				builder.Fields(tableFields...)
			} else {
				for _, field := range test.fields {
					mockColumn := mock.NewMockColumn(ctrl)
					mockColumn.
						EXPECT().
						SQLColumnName().
						Return(field)

					fields = append(fields, mockColumn)
				}

				table.
					EXPECT().
					Columns().
					Return(fields)
			}

			values := make([]*mock.MockBasicTable, 0, len(test.values))
			for _, value := range test.values {
				table := mock.NewMockBasicTable(ctrl)
				table.
					EXPECT().
					ColumnMap().
					Return(value)

				values = append(values, table)
			}
			builder = builder.Values(values...)

			query, args, err := builder.BuildQuery()

			if test.err {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			assert.Equal(t, test.query, query)
			assert.Equal(t, test.args, args)
		})
	}
}
