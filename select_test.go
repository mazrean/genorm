package genorm_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestSelectBuildQuery(t *testing.T) {
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

	tests := []struct {
		description     string
		tableExpr       expr
		distinct        bool
		isFieldSet      bool
		fields          []field
		groupExprs      []expr
		havingCondition *expr
		whereCondition  *expr
		orderItems      []orderItem
		limit           uint64
		offset          uint64
		lockType        genorm.LockType
		query           string
		args            []genorm.ExprType
		err             bool
	}{
		{
			description: "normal",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge`",
			args:  []genorm.ExprType{},
		},
		{
			description: "table error",
			tableExpr: expr{
				errs: []error{errors.New("table error")},
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			err: true,
		},
		{
			description: "multiple fields",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
				{
					tableName:     "hoge",
					columnName:    "piyo",
					sqlColumnName: "`hoge`.`piyo`",
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0, `hoge`.`piyo` AS hoge_piyo_0 FROM `hoge`",
			args:  []genorm.ExprType{},
		},
		{
			description: "field set",
			tableExpr: expr{
				query: "`hoge`",
			},
			isFieldSet: true,
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge`",
			args:  []genorm.ExprType{},
		},
		{
			description: "multiple fields set",
			tableExpr: expr{
				query: "`hoge`",
			},
			isFieldSet: true,
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
				{
					tableName:     "hoge",
					columnName:    "piyo",
					sqlColumnName: "`hoge`.`piyo`",
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0, `hoge`.`piyo` AS hoge_piyo_0 FROM `hoge`",
			args:  []genorm.ExprType{},
		},
		{
			description: "joined table",
			tableExpr: expr{
				query: "`hoge` JOIN `fuga` ON `hoge`.`id` = `fuga`.`id` AND `hoge`.`huga` = ?",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` JOIN `fuga` ON `hoge`.`id` = `fuga`.`id` AND `hoge`.`huga` = ?",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "distinct",
			tableExpr: expr{
				query: "`hoge`",
			},
			distinct: true,
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			query: "SELECT DISTINCT `hoge`.`huga` AS hoge_huga_0 FROM `hoge`",
			args:  []genorm.ExprType{},
		},
		{
			description: "group by",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			groupExprs: []expr{
				{
					query: "`hoge`.`fuga`",
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` GROUP BY `hoge`.`fuga`",
			args:  []genorm.ExprType{},
		},
		{
			description: "group by with args",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			groupExprs: []expr{
				{
					query: "`hoge`.`fuga` = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` GROUP BY `hoge`.`fuga` = ?",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "group by error",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			groupExprs: []expr{
				{
					errs: []error{errors.New("group error")},
				},
			},
			err: true,
		},
		{
			description: "multiple group by",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			groupExprs: []expr{
				{
					query: "`hoge`.`fuga` = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
				{
					query: "`hoge`.`piyo` = ?",
					args:  []genorm.ExprType{genorm.Wrap(2)},
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` GROUP BY `hoge`.`fuga` = ?, `hoge`.`piyo` = ?",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "having",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			groupExprs: []expr{
				{
					query: "`hoge`.`fuga`",
				},
			},
			havingCondition: &expr{
				query: "`hoge`.`huga` = ?",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` GROUP BY `hoge`.`fuga` HAVING `hoge`.`huga` = ?",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "having error",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			groupExprs: []expr{
				{
					query: "`hoge`.`fuga`",
				},
			},
			havingCondition: &expr{
				errs: []error{errors.New("having error")},
			},
			err: true,
		},
		{
			description: "where",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			whereCondition: &expr{
				query: "(`hoge`.`huga` = ?)",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` WHERE (`hoge`.`huga` = ?)",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "where error",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			whereCondition: &expr{
				errs: []error{errors.New("where error")},
			},
			err: true,
		},
		{
			description: "order by",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			orderItems: []orderItem{
				{
					direction: genorm.Asc,
					expr: expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` ORDER BY (`hoge`.`huga` = ?) ASC",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "order by error",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			orderItems: []orderItem{
				{
					direction: genorm.Asc,
					expr: expr{
						errs: []error{errors.New("order by error")},
					},
				},
			},
			err: true,
		},
		{
			description: "multi order by",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			orderItems: []orderItem{
				{
					direction: genorm.Asc,
					expr: expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
				},
				{
					direction: genorm.Desc,
					expr: expr{
						query: "(`hoge`.`nya` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(2)},
					},
				},
			},
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` ORDER BY (`hoge`.`huga` = ?) ASC, (`hoge`.`nya` = ?) DESC",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "limit",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			limit: 1,
			query: "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` LIMIT 1",
			args:  []genorm.ExprType{},
		},
		{
			description: "offset",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			offset: 1,
			query:  "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` OFFSET 1",
			args:   []genorm.ExprType{},
		},
		{
			description: "for update",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			lockType: genorm.ForUpdate,
			query:    "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` FOR UPDATE",
			args:     []genorm.ExprType{},
		},
		{
			description: "for share",
			tableExpr: expr{
				query: "`hoge`",
			},
			fields: []field{
				{
					tableName:     "hoge",
					columnName:    "huga",
					sqlColumnName: "`hoge`.`huga`",
				},
			},
			lockType: genorm.ForShare,
			query:    "SELECT `hoge`.`huga` AS hoge_huga_0 FROM `hoge` FOR SHARE",
			args:     []genorm.ExprType{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			table := mock.NewMockTable(ctrl)
			table.
				EXPECT().
				Expr().
				Return(test.tableExpr.query, test.tableExpr.args, test.tableExpr.errs)
			table.
				EXPECT().
				GetErrors().
				Return(nil)

			builder := genorm.Select(table)

			fields := make([]genorm.Column, 0, len(test.fields))
			if test.isFieldSet {
				tableFields := make([]genorm.TableColumns[*mock.MockTable], 0, len(test.fields))
				for _, field := range test.fields {
					mockColumn := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
					mockColumn.
						EXPECT().
						TableName().
						Return(field.tableName)
					mockColumn.
						EXPECT().
						ColumnName().
						Return(field.columnName)
					mockColumn.
						EXPECT().
						SQLColumnName().
						Return(field.sqlColumnName)

					tableFields = append(tableFields, mockColumn)
					fields = append(fields, mockColumn)
				}

				builder.Fields(tableFields...)
			} else {
				for _, field := range test.fields {
					mockColumn := mock.NewMockColumn(ctrl)
					mockColumn.
						EXPECT().
						TableName().
						Return(field.tableName)
					mockColumn.
						EXPECT().
						ColumnName().
						Return(field.columnName)
					mockColumn.
						EXPECT().
						SQLColumnName().
						Return(field.sqlColumnName)

					fields = append(fields, mockColumn)
				}

				table.
					EXPECT().
					Columns().
					Return(fields)
			}

			if test.distinct {
				builder = builder.Distinct()
			}

			if test.groupExprs != nil {
				groupExprs := make([]genorm.TableExpr[*mock.MockTable], 0, len(test.groupExprs))
				for _, groupExpr := range test.groupExprs {
					mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
					mockExpr.
						EXPECT().
						Expr().
						Return(groupExpr.query, groupExpr.args, groupExpr.errs)

					groupExprs = append(groupExprs, mockExpr)
				}

				builder = builder.GroupBy(groupExprs...)
			}

			if test.havingCondition != nil {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				mockExpr.
					EXPECT().
					Expr().
					Return(test.havingCondition.query, test.havingCondition.args, test.havingCondition.errs)

				builder = builder.Having(mockExpr)
			}

			if test.whereCondition != nil {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				mockExpr.
					EXPECT().
					Expr().
					Return(test.whereCondition.query, test.whereCondition.args, test.whereCondition.errs)

				builder = builder.Where(mockExpr)
			}

			for _, orderItem := range test.orderItems {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				mockExpr.
					EXPECT().
					Expr().
					Return(orderItem.expr.query, orderItem.expr.args, orderItem.expr.errs)
				builder = builder.OrderBy(orderItem.direction, mockExpr)
			}

			if test.limit > 0 {
				builder = builder.Limit(test.limit)
			}

			if test.offset > 0 {
				builder = builder.Offset(test.offset)
			}

			if test.lockType != 0 {
				builder = builder.Lock(test.lockType)
			}

			columns, query, args, err := builder.BuildQuery()

			if test.err {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			assert.Equal(t, fields, columns)
			assert.Equal(t, test.query, query)
			assert.Equal(t, test.args, args)
		})
	}
}
