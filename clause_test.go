package genorm_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestWhereConditionClauseSetTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		beforeExpr  bool
		setExpr     bool
		err         bool
	}{
		{
			description: "normal",
			setExpr:     true,
		},
		{
			description: "condition already set",
			beforeExpr:  true,
			setExpr:     true,
			err:         true,
		},
		{
			description: "nil expr",
			err:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var beforeExpr genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if !test.beforeExpr {
				beforeExpr = nil
			} else {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				beforeExpr = mockExpr
			}

			var setExpr genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if !test.setExpr {
				setExpr = nil
			} else {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				setExpr = mockExpr
			}

			c := genorm.NewWhereConditionClause(beforeExpr)

			err := c.Set(setExpr)

			if test.err {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			assert.Equal(t, setExpr, c.GetCondition())
		})
	}
}

func TestWhereConditionClauseExistTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		exist       bool
	}{
		{
			description: "normal",
			exist:       true,
		},
		{
			description: "not exist",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if !test.exist {
				expr = nil
			} else {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr = mockExpr
			}

			c := genorm.NewWhereConditionClause(expr)

			res := c.Exists()

			assert.Equal(t, test.exist, res)
		})
	}
}

func TestWhereConditionClauseGetExprTest(t *testing.T) {
	t.Parallel()

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	tests := []struct {
		description string
		expr        *expr
		err         bool
	}{
		{
			description: "normal",
			expr: &expr{
				query: "(`hoge`.`huga` = ?)",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
		},
		{
			description: "empty condition",
			err:         true,
		},
		{
			description: "error condition",
			expr: &expr{
				errs: []error{errors.New("error")},
			},
			err: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var expr genorm.TypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]]
			if test.expr == nil {
				expr = nil
			} else {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr = mockExpr

				mockExpr.
					EXPECT().
					Expr().
					Return(test.expr.query, test.expr.args, test.expr.errs)
			}

			c := genorm.NewWhereConditionClause(expr)

			query, args, err := c.GetExpr()

			if test.err {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			assert.Equal(t, test.expr.query, query)
			assert.Equal(t, test.expr.args, args)
		})
	}
}

func TestGroupClauseSetTest(t *testing.T) {
	t.Parallel()

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	tests := []struct {
		description string
		beforeExprs []*expr
		setExprs    []*expr
		err         bool
	}{
		{
			description: "normal",
			setExprs: []*expr{
				{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
		},
		{
			description: "empty before expr",
			beforeExprs: []*expr{},
			setExprs: []*expr{
				{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
		},
		{
			description: "condition already set",
			beforeExprs: []*expr{
				{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			setExprs: []*expr{
				{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(2)},
				},
			},
			err: true,
		},
		{
			description: "nil expr",
			err:         true,
		},
		{
			description: "no expr",
			setExprs:    []*expr{},
			err:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var beforeExprs []genorm.TableExpr[*mock.MockTable]
			if test.beforeExprs == nil {
				beforeExprs = nil
			} else {
				beforeExprs = make([]genorm.TableExpr[*mock.MockTable], 0, len(test.beforeExprs))
				for range test.beforeExprs {
					beforeExprs = append(beforeExprs, mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl))
				}
			}

			var setExprs []genorm.TableExpr[*mock.MockTable]
			if test.setExprs == nil {
				setExprs = nil
			} else {
				setExprs = make([]genorm.TableExpr[*mock.MockTable], 0, len(test.beforeExprs))
				for range test.setExprs {
					setExprs = append(setExprs, mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl))
				}
			}

			c := genorm.NewGroupClause(beforeExprs)

			err := c.Set(setExprs)

			if test.err {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			assert.Equal(t, setExprs, c.GetCondition())
		})
	}
}

func TestGroupClauseExistTest(t *testing.T) {
	t.Parallel()

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	tests := []struct {
		description string
		exprs       []*expr
		exist       bool
	}{
		{
			description: "normal",
			exprs: []*expr{
				{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			exist: true,
		},
		{
			description: "nil expr",
		},
		{
			description: "no expr",
			exprs:       []*expr{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var exprs []genorm.TableExpr[*mock.MockTable]
			if test.exprs == nil {
				exprs = nil
			} else {
				exprs = make([]genorm.TableExpr[*mock.MockTable], 0, len(test.exprs))
				for range test.exprs {
					exprs = append(exprs, mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl))
				}
			}

			c := genorm.NewGroupClause(exprs)

			res := c.Exists()

			assert.Equal(t, test.exist, res)
		})
	}
}

func TestGroupClauseGetExprTest(t *testing.T) {
	t.Parallel()

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	tests := []struct {
		description string
		exprs       []*expr
		query       string
		args        []genorm.ExprType
		err         bool
	}{
		{
			description: "normal",
			exprs: []*expr{
				{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			query: "GROUP BY (`hoge`.`huga` = ?)",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "multi exprs",
			exprs: []*expr{
				{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
				{
					query: "(`hoge`.`nya` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(2)},
				},
			},
			query: "GROUP BY (`hoge`.`huga` = ?), (`hoge`.`nya` = ?)",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "nil expr",
			err:         true,
		},
		{
			description: "no expr",
			exprs:       []*expr{},
			err:         true,
		},
		{
			description: "error condition",
			exprs: []*expr{
				{
					errs: []error{errors.New("error")},
				},
			},
			err: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var exprs []genorm.TableExpr[*mock.MockTable]
			if test.exprs == nil {
				exprs = nil
			} else {
				exprs = make([]genorm.TableExpr[*mock.MockTable], 0, len(test.exprs))
				for _, expr := range test.exprs {
					mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
					exprs = append(exprs, mockExpr)

					mockExpr.
						EXPECT().
						Expr().
						Return(expr.query, expr.args, expr.errs)
				}
			}

			c := genorm.NewGroupClause(exprs)

			query, args, err := c.GetExpr()

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
