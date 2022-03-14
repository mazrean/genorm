package genorm_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBuildQuery(t *testing.T) {
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
		tableName      string
		whereCondition *expr
		orderItems     []orderItem
		limit          uint64
		query          string
		args           []genorm.ExprType
		err            bool
	}{
		{
			description: "normal",
			tableName:   "hoge",
			query:       "DELETE FROM `hoge`",
			args:        []genorm.ExprType{},
		},
		{
			description: "where",
			tableName:   "hoge",
			whereCondition: &expr{
				query: "(`hoge`.`huga` = ?)",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			query: "DELETE FROM `hoge` WHERE (`hoge`.`huga` = ?)",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "where error",
			tableName:   "hoge",
			whereCondition: &expr{
				errs: []error{errors.New("where error")},
			},
			err: true,
		},
		{
			description: "order by",
			tableName:   "hoge",
			orderItems: []orderItem{
				{
					direction: genorm.Asc,
					expr: expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
				},
			},
			query: "DELETE FROM `hoge` ORDER BY (`hoge`.`huga` = ?) ASC",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "order by error",
			tableName:   "hoge",
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
			tableName:   "hoge",
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
			query: "DELETE FROM `hoge` ORDER BY (`hoge`.`huga` = ?) ASC, (`hoge`.`nya` = ?) DESC",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "limit",
			tableName:   "hoge",
			limit:       1,
			query:       "DELETE FROM `hoge` LIMIT 1",
			args:        []genorm.ExprType{},
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

			builder := genorm.Delete(table)

			if test.whereCondition != nil {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockBasicTable, genorm.WrappedPrimitive[bool]](ctrl)
				mockExpr.
					EXPECT().
					Expr().
					Return(test.whereCondition.query, test.whereCondition.args, test.whereCondition.errs)

				builder = builder.Where(mockExpr)
			}

			for _, orderItem := range test.orderItems {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockBasicTable, genorm.WrappedPrimitive[bool]](ctrl)
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
