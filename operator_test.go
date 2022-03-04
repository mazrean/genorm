package genorm_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		expr1IsNil  bool
		expr1Query  string
		expr1Args   []genorm.ExprType
		expr1Errs   []error
		expr2IsNil  bool
		expr2Query  string
		expr2Args   []genorm.ExprType
		expr2Errs   []error
		expected    *genorm.TableAssignExpr[*mock.MockTable]
		isError     bool
	}{
		{
			description: "normal",
			expr1Query:  "`hoge`.`huga`",
			expr1Args:   nil,
			expr2Query:  "(`hoge`.`huga` + ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(1)},
			expected: genorm.NewTableAssignExpr[*mock.MockTable](
				"`hoge`.`huga` = (`hoge`.`huga` + ?)",
				[]genorm.ExprType{genorm.Wrap(1)},
				nil,
			),
		},
		{
			description: "nil expr1",
			expr1IsNil:  true,
			expr2Query:  "(`hoge`.`huga` + ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(1)},
			isError:     true,
		},
		{
			description: "nil expr2",
			expr1Query:  "`hoge`.`huga`",
			expr1Args:   nil,
			expr2IsNil:  true,
			isError:     true,
		},
		{
			description: "expr1 error",
			expr1Errs:   []error{errors.New("expr1 error")},
			expr2Query:  "(`hoge`.`huga` + ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(1)},
			isError:     true,
		},
		{
			description: "expr2 error",
			expr1Query:  "`hoge`.`huga`",
			expr1Args:   nil,
			expr2Errs:   []error{errors.New("expr2 error")},
			isError:     true,
		},
		{
			description: "expr1 with args",
			expr1Query:  "`hoge`.`huga`",
			expr1Args:   []genorm.ExprType{genorm.Wrap(2)},
			expr2Query:  "(`hoge`.`huga` + ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(1)},
			expected: genorm.NewTableAssignExpr[*mock.MockTable](
				"`hoge`.`huga` = (`hoge`.`huga` + ?)",
				[]genorm.ExprType{genorm.Wrap(2), genorm.Wrap(1)},
				nil,
			),
		},
		{
			description: "expr2 no args",
			expr1Query:  "`hoge`.`huga`",
			expr1Args:   nil,
			expr2Query:  "(`hoge`.`huga` + `hoge`.`huga`)",
			expr2Args:   nil,
			expected: genorm.NewTableAssignExpr[*mock.MockTable](
				"`hoge`.`huga` = (`hoge`.`huga` + `hoge`.`huga`)",
				nil,
				nil,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr1 genorm.TypedTableColumns[*mock.MockTable, genorm.WrappedPrimitive[int]]
			if test.expr1IsNil {
				expr1 = nil
			} else {
				mockExpr1 := mock.NewMockTypedTableColumn[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)
				expr1 = mockExpr1

				if !test.expr2IsNil {
					mockExpr1.
						EXPECT().
						Expr().
						Return(test.expr1Query, test.expr1Args, test.expr1Errs)
				}
			}

			var expr2 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[int]]
			if test.expr2IsNil {
				expr2 = nil
			} else {
				mockExpr2 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[int]](ctrl)
				expr2 = mockExpr2

				if !test.expr1IsNil {
					mockExpr2.
						EXPECT().
						Expr().
						Return(test.expr2Query, test.expr2Args, test.expr2Errs)
				}
			}

			res := genorm.Assign(expr1, expr2)

			assert.NotNil(t, res)

			query, args, errs := res.AssignExpr()
			if test.isError {
				assert.NotNil(t, errs)
				assert.NotEmpty(t, errs)

				return
			}

			expectQuery, expectArgs, _ := test.expected.AssignExpr()

			if !assert.Nil(t, errs) {
				return
			}

			assert.Equal(t, expectQuery, query)
			assert.Equal(t, expectArgs, args)
		})
	}
}

func TestAnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		expr1IsNil    bool
		expr1Query    string
		expr1Args     []genorm.ExprType
		expr1Errs     []error
		expr2IsNil    bool
		expr2Query    string
		expr2Args     []genorm.ExprType
		expr2Errs     []error
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description:   "normal",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = ?) AND (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "nil expr1",
			expr1IsNil:  true,
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "nil expr2",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2IsNil:  true,
			isError:     true,
		},
		{
			description: "expr1 error",
			expr1Errs:   []error{errors.New("expr1 error")},
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "expr2 error",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2Errs:   []error{errors.New("expr2 error")},
			isError:     true,
		},
		{
			description:   "expr1 no args",
			expr1Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr1Args:     nil,
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = `hoge`.`huga`) AND (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(2)},
		},
		{
			description:   "expr2 no args",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr2Args:     nil,
			expectedQuery: "((`hoge`.`huga` = ?) AND (`hoge`.`huga` = `hoge`.`huga`))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr1 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr1IsNil {
				expr1 = nil
			} else {
				mockExpr1 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr1 = mockExpr1

				if !test.expr2IsNil {
					mockExpr1.
						EXPECT().
						Expr().
						Return(test.expr1Query, test.expr1Args, test.expr1Errs)
				}
			}

			var expr2 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr2IsNil {
				expr2 = nil
			} else {
				mockExpr2 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr2 = mockExpr2

				if !test.expr1IsNil {
					mockExpr2.
						EXPECT().
						Expr().
						Return(test.expr2Query, test.expr2Args, test.expr2Errs)
				}
			}

			res := genorm.And(expr1, expr2)

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

func TestOr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		expr1IsNil    bool
		expr1Query    string
		expr1Args     []genorm.ExprType
		expr1Errs     []error
		expr2IsNil    bool
		expr2Query    string
		expr2Args     []genorm.ExprType
		expr2Errs     []error
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description:   "normal",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = ?) OR (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "nil expr1",
			expr1IsNil:  true,
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "nil expr2",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2IsNil:  true,
			isError:     true,
		},
		{
			description: "expr1 error",
			expr1Errs:   []error{errors.New("expr1 error")},
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "expr2 error",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2Errs:   []error{errors.New("expr2 error")},
			isError:     true,
		},
		{
			description:   "expr1 no args",
			expr1Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr1Args:     nil,
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = `hoge`.`huga`) OR (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(2)},
		},
		{
			description:   "expr2 no args",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr2Args:     nil,
			expectedQuery: "((`hoge`.`huga` = ?) OR (`hoge`.`huga` = `hoge`.`huga`))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr1 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr1IsNil {
				expr1 = nil
			} else {
				mockExpr1 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr1 = mockExpr1

				if !test.expr2IsNil {
					mockExpr1.
						EXPECT().
						Expr().
						Return(test.expr1Query, test.expr1Args, test.expr1Errs)
				}
			}

			var expr2 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr2IsNil {
				expr2 = nil
			} else {
				mockExpr2 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr2 = mockExpr2

				if !test.expr1IsNil {
					mockExpr2.
						EXPECT().
						Expr().
						Return(test.expr2Query, test.expr2Args, test.expr2Errs)
				}
			}

			res := genorm.Or(expr1, expr2)

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

func TestXor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		expr1IsNil    bool
		expr1Query    string
		expr1Args     []genorm.ExprType
		expr1Errs     []error
		expr2IsNil    bool
		expr2Query    string
		expr2Args     []genorm.ExprType
		expr2Errs     []error
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description:   "normal",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = ?) XOR (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "nil expr1",
			expr1IsNil:  true,
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "nil expr2",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2IsNil:  true,
			isError:     true,
		},
		{
			description: "expr1 error",
			expr1Errs:   []error{errors.New("expr1 error")},
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "expr2 error",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2Errs:   []error{errors.New("expr2 error")},
			isError:     true,
		},
		{
			description:   "expr1 no args",
			expr1Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr1Args:     nil,
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = `hoge`.`huga`) XOR (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(2)},
		},
		{
			description:   "expr2 no args",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr2Args:     nil,
			expectedQuery: "((`hoge`.`huga` = ?) XOR (`hoge`.`huga` = `hoge`.`huga`))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr1 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr1IsNil {
				expr1 = nil
			} else {
				mockExpr1 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr1 = mockExpr1

				if !test.expr2IsNil {
					mockExpr1.
						EXPECT().
						Expr().
						Return(test.expr1Query, test.expr1Args, test.expr1Errs)
				}
			}

			var expr2 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr2IsNil {
				expr2 = nil
			} else {
				mockExpr2 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr2 = mockExpr2

				if !test.expr1IsNil {
					mockExpr2.
						EXPECT().
						Expr().
						Return(test.expr2Query, test.expr2Args, test.expr2Errs)
				}
			}

			res := genorm.Xor(expr1, expr2)

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

func TestNot(t *testing.T) {
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
			expectedQuery: "(NOT (`hoge`.`huga` = ?))",
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
			expectedQuery: "(NOT (`hoge`.`huga` = `hoge`.`huga`))",
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

			res := genorm.Not(expr)

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

func TestEq(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		expr1IsNil    bool
		expr1Query    string
		expr1Args     []genorm.ExprType
		expr1Errs     []error
		expr2IsNil    bool
		expr2Query    string
		expr2Args     []genorm.ExprType
		expr2Errs     []error
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description:   "normal",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = ?) = (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "nil expr1",
			expr1IsNil:  true,
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "nil expr2",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2IsNil:  true,
			isError:     true,
		},
		{
			description: "expr1 error",
			expr1Errs:   []error{errors.New("expr1 error")},
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "expr2 error",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2Errs:   []error{errors.New("expr2 error")},
			isError:     true,
		},
		{
			description:   "expr1 no args",
			expr1Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr1Args:     nil,
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = `hoge`.`huga`) = (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(2)},
		},
		{
			description:   "expr2 no args",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr2Args:     nil,
			expectedQuery: "((`hoge`.`huga` = ?) = (`hoge`.`huga` = `hoge`.`huga`))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr1 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr1IsNil {
				expr1 = nil
			} else {
				mockExpr1 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr1 = mockExpr1

				if !test.expr2IsNil {
					mockExpr1.
						EXPECT().
						Expr().
						Return(test.expr1Query, test.expr1Args, test.expr1Errs)
				}
			}

			var expr2 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr2IsNil {
				expr2 = nil
			} else {
				mockExpr2 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr2 = mockExpr2

				if !test.expr1IsNil {
					mockExpr2.
						EXPECT().
						Expr().
						Return(test.expr2Query, test.expr2Args, test.expr2Errs)
				}
			}

			res := genorm.Eq(expr1, expr2)

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

func TestNeq(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		expr1IsNil    bool
		expr1Query    string
		expr1Args     []genorm.ExprType
		expr1Errs     []error
		expr2IsNil    bool
		expr2Query    string
		expr2Args     []genorm.ExprType
		expr2Errs     []error
		expectedQuery string
		expectedArgs  []genorm.ExprType
		isError       bool
	}{
		{
			description:   "normal",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = ?) != (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "nil expr1",
			expr1IsNil:  true,
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "nil expr2",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2IsNil:  true,
			isError:     true,
		},
		{
			description: "expr1 error",
			expr1Errs:   []error{errors.New("expr1 error")},
			expr2Query:  "(`hoge`.`huga` > ?)",
			expr2Args:   []genorm.ExprType{genorm.Wrap(2)},
			isError:     true,
		},
		{
			description: "expr2 error",
			expr1Query:  "(`hoge`.`huga` = ?)",
			expr1Args:   []genorm.ExprType{genorm.Wrap(1)},
			expr2Errs:   []error{errors.New("expr2 error")},
			isError:     true,
		},
		{
			description:   "expr1 no args",
			expr1Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr1Args:     nil,
			expr2Query:    "(`hoge`.`huga` > ?)",
			expr2Args:     []genorm.ExprType{genorm.Wrap(2)},
			expectedQuery: "((`hoge`.`huga` = `hoge`.`huga`) != (`hoge`.`huga` > ?))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(2)},
		},
		{
			description:   "expr2 no args",
			expr1Query:    "(`hoge`.`huga` = ?)",
			expr1Args:     []genorm.ExprType{genorm.Wrap(1)},
			expr2Query:    "(`hoge`.`huga` = `hoge`.`huga`)",
			expr2Args:     nil,
			expectedQuery: "((`hoge`.`huga` = ?) != (`hoge`.`huga` = `hoge`.`huga`))",
			expectedArgs:  []genorm.ExprType{genorm.Wrap(1)},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr1 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr1IsNil {
				expr1 = nil
			} else {
				mockExpr1 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr1 = mockExpr1

				if !test.expr2IsNil {
					mockExpr1.
						EXPECT().
						Expr().
						Return(test.expr1Query, test.expr1Args, test.expr1Errs)
				}
			}

			var expr2 genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr2IsNil {
				expr2 = nil
			} else {
				mockExpr2 := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr2 = mockExpr2

				if !test.expr1IsNil {
					mockExpr2.
						EXPECT().
						Expr().
						Return(test.expr2Query, test.expr2Args, test.expr2Errs)
				}
			}

			res := genorm.Neq(expr1, expr2)

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
