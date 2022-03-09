package genorm_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestAvg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		exprIsNil     bool
		exprQuery     string
		exprArgs      []genorm.ExprType
		exprErrs      []error
		distinct      bool
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description:   "normal",
			exprQuery:     "(`hoge`.`huga` = ?)",
			exprArgs:      []genorm.ExprType{genorm.Wrap(1)},
			expectedQuery: "AVG((`hoge`.`huga` = ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description:   "distinct",
			exprQuery:     "(`hoge`.`huga` = ?)",
			exprArgs:      []genorm.ExprType{genorm.Wrap(1)},
			distinct:      true,
			expectedQuery: "AVG(DISTINCT (`hoge`.`huga` = ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "nil expr",
			exprIsNil:   true,
			isError:     true,
		},
		{
			description: "expr error",
			exprErrs:    []error{errors.New("expr1 error")},
			isError:     true,
		},
		{
			description:   "expr1 no args",
			exprQuery:     "(`hoge`.`huga` = `hoge`.`huga`)",
			exprArgs:      nil,
			expectedQuery: "AVG((`hoge`.`huga` = `hoge`.`huga`))",
			expectedArgs:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.exprIsNil {
				expr = nil
			} else {
				mockExpr := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr = mockExpr

				mockExpr.
					EXPECT().
					Expr().
					Return(test.exprQuery, test.exprArgs, test.exprErrs)
			}

			res := genorm.Avg(expr, test.distinct)

			assert.NotNil(t, res)

			query, args, errs := res.Expr()
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

func TestCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		exprIsNil     bool
		exprQuery     string
		exprArgs      []genorm.ExprType
		exprErrs      []error
		distinct      bool
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description:   "normal",
			exprQuery:     "(`hoge`.`huga` = ?)",
			exprArgs:      []genorm.ExprType{genorm.Wrap(1)},
			expectedQuery: "COUNT((`hoge`.`huga` = ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description:   "distinct",
			exprQuery:     "(`hoge`.`huga` = ?)",
			exprArgs:      []genorm.ExprType{genorm.Wrap(1)},
			distinct:      true,
			expectedQuery: "COUNT(DISTINCT (`hoge`.`huga` = ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "nil expr",
			exprIsNil:   true,
			isError:     true,
		},
		{
			description: "expr error",
			exprErrs:    []error{errors.New("expr1 error")},
			isError:     true,
		},
		{
			description:   "expr1 no args",
			exprQuery:     "(`hoge`.`huga` = `hoge`.`huga`)",
			exprArgs:      nil,
			expectedQuery: "COUNT((`hoge`.`huga` = `hoge`.`huga`))",
			expectedArgs:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.exprIsNil {
				expr = nil
			} else {
				mockExpr := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr = mockExpr

				mockExpr.
					EXPECT().
					Expr().
					Return(test.exprQuery, test.exprArgs, test.exprErrs)
			}

			res := genorm.Count(expr, test.distinct)

			assert.NotNil(t, res)

			query, args, errs := res.Expr()
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

func TestMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		exprIsNil     bool
		exprQuery     string
		exprArgs      []genorm.ExprType
		exprErrs      []error
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description:   "normal",
			exprQuery:     "(`hoge`.`huga` = ?)",
			exprArgs:      []genorm.ExprType{genorm.Wrap(1)},
			expectedQuery: "MAX((`hoge`.`huga` = ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "nil expr",
			exprIsNil:   true,
			isError:     true,
		},
		{
			description: "expr error",
			exprErrs:    []error{errors.New("expr1 error")},
			isError:     true,
		},
		{
			description:   "expr1 no args",
			exprQuery:     "(`hoge`.`huga` = `hoge`.`huga`)",
			exprArgs:      nil,
			expectedQuery: "MAX((`hoge`.`huga` = `hoge`.`huga`))",
			expectedArgs:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.exprIsNil {
				expr = nil
			} else {
				mockExpr := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr = mockExpr

				mockExpr.
					EXPECT().
					Expr().
					Return(test.exprQuery, test.exprArgs, test.exprErrs)
			}

			res := genorm.Max(expr)

			assert.NotNil(t, res)

			query, args, errs := res.Expr()
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

func TestMin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		exprIsNil     bool
		exprQuery     string
		exprArgs      []genorm.ExprType
		exprErrs      []error
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description:   "normal",
			exprQuery:     "(`hoge`.`huga` = ?)",
			exprArgs:      []genorm.ExprType{genorm.Wrap(1)},
			expectedQuery: "MIN((`hoge`.`huga` = ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "nil expr",
			exprIsNil:   true,
			isError:     true,
		},
		{
			description: "expr error",
			exprErrs:    []error{errors.New("expr1 error")},
			isError:     true,
		},
		{
			description:   "expr1 no args",
			exprQuery:     "(`hoge`.`huga` = `hoge`.`huga`)",
			exprArgs:      nil,
			expectedQuery: "MIN((`hoge`.`huga` = `hoge`.`huga`))",
			expectedArgs:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.exprIsNil {
				expr = nil
			} else {
				mockExpr := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr = mockExpr

				mockExpr.
					EXPECT().
					Expr().
					Return(test.exprQuery, test.exprArgs, test.exprErrs)
			}

			res := genorm.Min(expr)

			assert.NotNil(t, res)

			query, args, errs := res.Expr()
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
