package genorm_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestSubQuery(t *testing.T) {
	t.Parallel()

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	tests := []struct {
		description   string
		pluckCtxIsNil bool
		tableExpr     expr
		fieldExpr     expr
		whereExpr     *expr
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description: "normal subquery",
			tableExpr: expr{
				query: "inner_table",
			},
			fieldExpr: expr{
				query: "inner_table.column1",
			},
			expectedQuery: "(SELECT inner_table.column1 AS res FROM inner_table)",
			expectedArgs:  []genorm.ExprType{},
		},
		{
			description: "subquery with where clause",
			tableExpr: expr{
				query: "inner_table",
			},
			fieldExpr: expr{
				query: "inner_table.column1",
			},
			whereExpr: &expr{
				query: "(inner_table.column2 = ?)",
				args:  []genorm.ExprType{genorm.Wrap(100)},
			},
			expectedQuery: "(SELECT inner_table.column1 AS res FROM inner_table WHERE (inner_table.column2 = ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(100)},
		},
		{
			description: "subquery with multiple args",
			tableExpr: expr{
				query: "inner_table JOIN other_table ON inner_table.id = other_table.inner_id",
			},
			fieldExpr: expr{
				query: "MAX(inner_table.value)",
			},
			whereExpr: &expr{
				query: "(inner_table.status = ? AND inner_table.active = ?)",
				args:  []genorm.ExprType{genorm.Wrap("active"), genorm.Wrap(true)},
			},
			expectedQuery: "(SELECT MAX(inner_table.value) AS res FROM inner_table JOIN other_table ON inner_table.id = other_table.inner_id WHERE (inner_table.status = ? AND inner_table.active = ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap("active"), genorm.Wrap(true)},
		},
		{
			description:   "nil pluck context",
			pluckCtxIsNil: true,
			isError:       true,
		},
		{
			description: "table error in pluck context",
			tableExpr: expr{
				errs: []error{errors.New("table error")},
			},
			fieldExpr: expr{
				query: "inner_table.column1",
			},
			isError: true,
		},
		{
			description: "field error in pluck context",
			tableExpr: expr{
				query: "inner_table",
			},
			fieldExpr: expr{
				errs: []error{errors.New("field error")},
			},
			isError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var pluckCtx *genorm.PluckContext[*mock.MockTable, genorm.WrappedPrimitive[int]]
			if test.pluckCtxIsNil {
				pluckCtx = nil
			} else {
				mockTable := mock.NewMockTable(ctrl)
				mockField := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)

				mockTable.
					EXPECT().
					GetErrors().
					Return(nil).
					AnyTimes()

				mockTable.
					EXPECT().
					Expr().
					Return(test.tableExpr.query, test.tableExpr.args, test.tableExpr.errs).
					AnyTimes()

				mockField.
					EXPECT().
					Expr().
					Return(test.fieldExpr.query, test.fieldExpr.args, test.fieldExpr.errs).
					AnyTimes()

				pluckCtx = genorm.Pluck(mockTable, mockField)

				if test.whereExpr != nil {
					mockWhere := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
					mockWhere.
						EXPECT().
						Expr().
						Return(test.whereExpr.query, test.whereExpr.args, test.whereExpr.errs).
						AnyTimes()

					pluckCtx = pluckCtx.Where(mockWhere)
				}
			}

			subQuery := genorm.SubQuery[*mock.MockTable](pluckCtx)

			assert.NotNil(t, subQuery)

			query, args, errs := subQuery.Expr()
			if test.isError {
				assert.NotNil(t, errs)
				assert.NotEmpty(t, errs)
				return
			}

			if !assert.Nil(t, errs) {
				return
			}

			assert.Equal(t, test.expectedQuery, query)
			assert.Equal(t, test.expectedArgs, args)
		})
	}
}

func TestSubQueryWithOperators(t *testing.T) {
	t.Parallel()

	t.Run("subquery with Eq operator", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		// Setup inner table for subquery
		mockInnerTable := mock.NewMockTable(ctrl)
		mockInnerField := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)

		mockInnerTable.
			EXPECT().
			GetErrors().
			Return(nil).
			AnyTimes()

		mockInnerTable.
			EXPECT().
			Expr().
			Return("inner_table", []genorm.ExprType{}, []error{}).
			AnyTimes()

		mockInnerField.
			EXPECT().
			Expr().
			Return("inner_table.max_value", []genorm.ExprType{}, []error{}).
			AnyTimes()

		// Setup outer table column
		mockOuterColumn := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)
		mockOuterColumn.
			EXPECT().
			Expr().
			Return("outer_table.value", []genorm.ExprType{}, []error{}).
			AnyTimes()

		// Create subquery
		pluckCtx := genorm.Pluck(mockInnerTable, mockInnerField)
		subQuery := genorm.SubQuery[*mock.MockTable](pluckCtx)

		// Use subquery in Eq operator
		result := genorm.Eq(mockOuterColumn, subQuery)

		query, args, errs := result.Expr()
		assert.Nil(t, errs)
		assert.Equal(t, "(outer_table.value = (SELECT inner_table.max_value AS res FROM inner_table))", query)
		assert.Equal(t, []genorm.ExprType{}, args)
	})

	t.Run("subquery with In operator", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		// Setup inner table for subquery
		mockInnerTable := mock.NewMockTable(ctrl)
		mockInnerField := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)

		mockInnerTable.
			EXPECT().
			GetErrors().
			Return(nil).
			AnyTimes()

		mockInnerTable.
			EXPECT().
			Expr().
			Return("inner_table", []genorm.ExprType{}, []error{}).
			AnyTimes()

		mockInnerField.
			EXPECT().
			Expr().
			Return("inner_table.id", []genorm.ExprType{}, []error{}).
			AnyTimes()

		// Setup outer table column
		mockOuterColumn := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)
		mockOuterColumn.
			EXPECT().
			Expr().
			Return("outer_table.id", []genorm.ExprType{}, []error{}).
			AnyTimes()

		// Create subquery
		pluckCtx := genorm.Pluck(mockInnerTable, mockInnerField)
		subQuery := genorm.SubQuery[*mock.MockTable](pluckCtx)

		// Use subquery in comparison
		result := genorm.Gt(mockOuterColumn, subQuery)

		query, args, errs := result.Expr()
		assert.Nil(t, errs)
		assert.Equal(t, "(outer_table.id > (SELECT inner_table.id AS res FROM inner_table))", query)
		assert.Equal(t, []genorm.ExprType{}, args)
	})
}
