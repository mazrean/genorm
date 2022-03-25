package genorm_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestFindBuildQuery(t *testing.T) {
	t.Parallel()

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
		isFieldErr      bool
		tupleExprs      []expr
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
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			query: "SELECT hoge.huga AS value0 FROM hoge",
			args:  []genorm.ExprType{},
		},
		{
			description: "table error",
			tableExpr: expr{
				errs: []error{errors.New("table error")},
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			err: true,
		},
		{
			description: "joined table",
			tableExpr: expr{
				query: "hoge JOIN fuga ON hoge.id = fuga.id AND hoge.huga = ?",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			query: "SELECT hoge.huga AS value0 FROM hoge JOIN fuga ON hoge.id = fuga.id AND hoge.huga = ?",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "field with args",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "(hoge.huga = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			query: "SELECT (hoge.huga = ?) AS value0 FROM hoge",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "multiple field",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
				{
					query: "hoge.piyo",
				},
			},
			query: "SELECT hoge.huga AS value0, hoge.piyo AS value1 FROM hoge",
			args:  []genorm.ExprType{},
		},
		{
			description: "field error",
			tableExpr: expr{
				query: "hoge",
			},
			isFieldErr: true,
			tupleExprs: []expr{
				{
					errs: []error{errors.New("field error")},
				},
			},
			err: true,
		},
		{
			description: "distinct",
			tableExpr: expr{
				query: "hoge",
			},
			distinct: true,
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			query: "SELECT DISTINCT hoge.huga AS value0 FROM hoge",
			args:  []genorm.ExprType{},
		},
		{
			description: "group by",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			groupExprs: []expr{
				{
					query: "hoge.fuga",
				},
			},
			query: "SELECT hoge.huga AS value0 FROM hoge GROUP BY hoge.fuga",
			args:  []genorm.ExprType{},
		},
		{
			description: "group by with args",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			groupExprs: []expr{
				{
					query: "hoge.fuga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			query: "SELECT hoge.huga AS value0 FROM hoge GROUP BY hoge.fuga = ?",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "group by error",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
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
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			groupExprs: []expr{
				{
					query: "hoge.fuga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
				{
					query: "hoge.piyo = ?",
					args:  []genorm.ExprType{genorm.Wrap(2)},
				},
			},
			query: "SELECT hoge.huga AS value0 FROM hoge GROUP BY hoge.fuga = ?, hoge.piyo = ?",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "having",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			groupExprs: []expr{
				{
					query: "hoge.fuga",
				},
			},
			havingCondition: &expr{
				query: "hoge.huga = ?",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			query: "SELECT hoge.huga AS value0 FROM hoge GROUP BY hoge.fuga HAVING hoge.huga = ?",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "having error",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			groupExprs: []expr{
				{
					query: "hoge.fuga",
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
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			whereCondition: &expr{
				query: "(hoge.huga = ?)",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			query: "SELECT hoge.huga AS value0 FROM hoge WHERE (hoge.huga = ?)",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "where error",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
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
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			orderItems: []orderItem{
				{
					direction: genorm.Asc,
					expr: expr{
						query: "(hoge.huga = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
				},
			},
			query: "SELECT hoge.huga AS value0 FROM hoge ORDER BY (hoge.huga = ?) ASC",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "order by error",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
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
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			orderItems: []orderItem{
				{
					direction: genorm.Asc,
					expr: expr{
						query: "(hoge.huga = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
				},
				{
					direction: genorm.Desc,
					expr: expr{
						query: "(hoge.nya = ?)",
						args:  []genorm.ExprType{genorm.Wrap(2)},
					},
				},
			},
			query: "SELECT hoge.huga AS value0 FROM hoge ORDER BY (hoge.huga = ?) ASC, (hoge.nya = ?) DESC",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "limit",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			limit: 1,
			query: "SELECT hoge.huga AS value0 FROM hoge LIMIT 1",
			args:  []genorm.ExprType{},
		},
		{
			description: "offset",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			offset: 1,
			query:  "SELECT hoge.huga AS value0 FROM hoge OFFSET 1",
			args:   []genorm.ExprType{},
		},
		{
			description: "for update",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			lockType: genorm.ForUpdate,
			query:    "SELECT hoge.huga AS value0 FROM hoge FOR UPDATE",
			args:     []genorm.ExprType{},
		},
		{
			description: "for share",
			tableExpr: expr{
				query: "hoge",
			},
			tupleExprs: []expr{
				{
					query: "hoge.huga",
				},
			},
			lockType: genorm.ForShare,
			query:    "SELECT hoge.huga AS value0 FROM hoge FOR SHARE",
			args:     []genorm.ExprType{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			table := mock.NewMockTable(ctrl)
			if !test.isFieldErr {
				table.
					EXPECT().
					Expr().
					Return(test.tableExpr.query, test.tableExpr.args, test.tableExpr.errs)
			}
			table.
				EXPECT().
				GetErrors().
				Return(nil)

			exprs := make([]genorm.Expr, 0, len(test.tupleExprs))
			for _, fieldExpr := range test.tupleExprs {
				mockField := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)
				mockField.
					EXPECT().
					Expr().
					Return(fieldExpr.query, fieldExpr.args, fieldExpr.errs)
				var field genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[int]] = mockField

				exprs = append(exprs, field)
			}
			mockTuple := mock.NewMockTuple(ctrl)
			mockTuple.
				EXPECT().
				Exprs().
				Return(exprs)

			builder := genorm.Find(table, mockTuple)

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
