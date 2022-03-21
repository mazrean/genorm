package genorm_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBuildQuery(t *testing.T) {
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
		description    string
		tableExpr      expr
		assignExprs    []expr
		whereCondition *expr
		orderItems     []orderItem
		limit          uint64
		query          string
		args           []genorm.ExprType
		err            bool
	}{
		{
			description: "normal",
			tableExpr: expr{
				query: "hoge",
			},
			assignExprs: []expr{
				{
					query: "hoge.huga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			query: "UPDATE hoge SET hoge.huga = ?",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "joined table",
			tableExpr: expr{
				query: "hoge JOIN fuga ON hoge.id = fuga.id AND hoge.huga = ?",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			assignExprs: []expr{
				{
					query: "hoge.huga = ?",
					args:  []genorm.ExprType{genorm.Wrap(2)},
				},
			},
			query: "UPDATE hoge JOIN fuga ON hoge.id = fuga.id AND hoge.huga = ? SET hoge.huga = ?",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "multi assign",
			tableExpr: expr{
				query: "hoge",
			},
			assignExprs: []expr{
				{
					query: "hoge.huga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
				{
					query: "hoge.nya = ?",
					args:  []genorm.ExprType{genorm.Wrap(2)},
				},
			},
			query: "UPDATE hoge SET hoge.huga = ?, hoge.nya = ?",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "where",
			tableExpr: expr{
				query: "hoge",
			},
			assignExprs: []expr{
				{
					query: "hoge.huga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			whereCondition: &expr{
				query: "(hoge.huga = ?)",
				args:  []genorm.ExprType{genorm.Wrap(2)},
			},
			query: "UPDATE hoge SET hoge.huga = ? WHERE (hoge.huga = ?)",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "where error",
			tableExpr: expr{
				query: "hoge",
			},
			assignExprs: []expr{
				{
					query: "hoge.huga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
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
			assignExprs: []expr{
				{
					query: "hoge.huga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			orderItems: []orderItem{
				{
					direction: genorm.Asc,
					expr: expr{
						query: "(hoge.huga = ?)",
						args:  []genorm.ExprType{genorm.Wrap(2)},
					},
				},
			},
			query: "UPDATE hoge SET hoge.huga = ? ORDER BY (hoge.huga = ?) ASC",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "order by error",
			tableExpr: expr{
				query: "hoge",
			},
			assignExprs: []expr{
				{
					query: "hoge.huga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
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
			assignExprs: []expr{
				{
					query: "hoge.huga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			orderItems: []orderItem{
				{
					direction: genorm.Asc,
					expr: expr{
						query: "(hoge.huga = ?)",
						args:  []genorm.ExprType{genorm.Wrap(2)},
					},
				},
				{
					direction: genorm.Desc,
					expr: expr{
						query: "(hoge.nya = ?)",
						args:  []genorm.ExprType{genorm.Wrap(3)},
					},
				},
			},
			query: "UPDATE hoge SET hoge.huga = ? ORDER BY (hoge.huga = ?) ASC, (hoge.nya = ?) DESC",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2), genorm.Wrap(3)},
		},
		{
			description: "limit",
			tableExpr: expr{
				query: "hoge",
			},
			assignExprs: []expr{
				{
					query: "hoge.huga = ?",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			limit: 1,
			query: "UPDATE hoge SET hoge.huga = ? LIMIT 1",
			args:  []genorm.ExprType{genorm.Wrap(1)},
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

			builder := genorm.Update(table)

			assignExprs := make([]*genorm.TableAssignExpr[*mock.MockTable], 0, len(test.assignExprs))
			for _, assignExpr := range test.assignExprs {
				assignExprs = append(assignExprs, genorm.NewTableAssignExpr[*mock.MockTable](assignExpr.query, assignExpr.args, assignExpr.errs))
			}
			builder = builder.Set(assignExprs...)

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
