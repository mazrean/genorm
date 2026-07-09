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

// TestSubQueryExampleUsage demonstrates typical usage patterns for SubQuery
func TestSubQueryExampleUsage(t *testing.T) {
	t.Parallel()

	t.Run("find users with value greater than subquery max", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		// Inner table for subquery - finds max value from statistics
		mockStatsTable := mock.NewMockTable(ctrl)
		mockStatsMaxValue := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)
		mockStatsStatusColumn := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)

		mockStatsTable.
			EXPECT().
			GetErrors().
			Return(nil).
			AnyTimes()

		mockStatsTable.
			EXPECT().
			Expr().
			Return("statistics", []genorm.ExprType{}, []error{}).
			AnyTimes()

		mockStatsMaxValue.
			EXPECT().
			Expr().
			Return("MAX(statistics.value)", []genorm.ExprType{}, []error{}).
			AnyTimes()

		mockStatsStatusColumn.
			EXPECT().
			Expr().
			Return("(statistics.active = ?)", []genorm.ExprType{genorm.Wrap(true)}, []error{}).
			AnyTimes()

		// Outer table column - user value
		mockUserValue := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)
		mockUserValue.
			EXPECT().
			Expr().
			Return("users.value", []genorm.ExprType{}, []error{}).
			AnyTimes()

		// Create subquery: SELECT MAX(value) FROM statistics WHERE active = true
		subQuery := genorm.SubQuery[*mock.MockTable](
			genorm.Pluck(mockStatsTable, mockStatsMaxValue).
				Where(mockStatsStatusColumn),
		)

		// Use in comparison: users.value > (SELECT MAX(value) FROM statistics WHERE active = true)
		condition := genorm.Gt(mockUserValue, subQuery)

		query, args, errs := condition.Expr()
		assert.Nil(t, errs)
		assert.Equal(t, "(users.value > (SELECT MAX(statistics.value) AS res FROM statistics WHERE (statistics.active = ?)))", query)
		assert.Equal(t, []genorm.ExprType{genorm.Wrap(true)}, args)
	})

	t.Run("find records where column equals subquery result", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		// Inner table for subquery
		mockInnerTable := mock.NewMockTable(ctrl)
		mockInnerColumn := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)

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

		mockInnerColumn.
			EXPECT().
			Expr().
			Return("inner_table.target_id", []genorm.ExprType{}, []error{}).
			AnyTimes()

		// Outer table column
		mockOuterColumn := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)
		mockOuterColumn.
			EXPECT().
			Expr().
			Return("outer_table.id", []genorm.ExprType{}, []error{}).
			AnyTimes()

		// Create subquery with LIMIT to get single value
		pluckCtx := genorm.Pluck(mockInnerTable, mockInnerColumn).Limit(1)
		subQuery := genorm.SubQuery[*mock.MockTable](pluckCtx)

		// Use in equality comparison
		condition := genorm.Eq(mockOuterColumn, subQuery)

		query, args, errs := condition.Expr()
		assert.Nil(t, errs)
		assert.Equal(t, "(outer_table.id = (SELECT inner_table.target_id AS res FROM inner_table LIMIT 1))", query)
		assert.Equal(t, []genorm.ExprType{}, args)
	})

	t.Run("subquery with count aggregate", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		// Inner table for subquery
		mockInnerTable := mock.NewMockTable(ctrl)
		mockInnerIDColumn := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)
		mockInnerStatusColumn := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)

		mockInnerTable.
			EXPECT().
			GetErrors().
			Return(nil).
			AnyTimes()

		mockInnerTable.
			EXPECT().
			Expr().
			Return("orders", []genorm.ExprType{}, []error{}).
			AnyTimes()

		mockInnerIDColumn.
			EXPECT().
			Expr().
			Return("orders.id", []genorm.ExprType{}, []error{}).
			AnyTimes()

		mockInnerStatusColumn.
			EXPECT().
			Expr().
			Return("(orders.status = ?)", []genorm.ExprType{genorm.Wrap("completed")}, []error{}).
			AnyTimes()

		// Outer table column
		mockOuterThreshold := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int64]](ctrl)
		mockOuterThreshold.
			EXPECT().
			Expr().
			Return("users.order_threshold", []genorm.ExprType{}, []error{}).
			AnyTimes()

		// Create COUNT subquery
		mockCountExpr := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int64]](ctrl)
		mockCountExpr.
			EXPECT().
			Expr().
			Return("COUNT(orders.id)", []genorm.ExprType{}, []error{}).
			AnyTimes()

		subQuery := genorm.SubQuery[*mock.MockTable](
			genorm.Pluck(mockInnerTable, mockCountExpr).
				Where(mockInnerStatusColumn),
		)

		// Use in comparison: threshold < (SELECT COUNT(id) FROM orders WHERE status = 'completed')
		condition := genorm.Lt(mockOuterThreshold, subQuery)

		query, args, errs := condition.Expr()
		assert.Nil(t, errs)
		assert.Equal(t, "(users.order_threshold < (SELECT COUNT(orders.id) AS res FROM orders WHERE (orders.status = ?)))", query)
		assert.Equal(t, []genorm.ExprType{genorm.Wrap("completed")}, args)
	})
}
