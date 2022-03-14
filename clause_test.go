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

func TestOrderClauseAddTest(t *testing.T) {
	t.Parallel()

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	type orderItem struct {
		expr      *expr
		direction genorm.OrderDirection
	}

	tests := []struct {
		description string
		beforeExprs []orderItem
		addItem     orderItem
		afterExprs  []orderItem
		err         bool
	}{
		{
			description: "normal",
			addItem: orderItem{
				expr: &expr{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
				direction: genorm.Asc,
			},
			afterExprs: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
					direction: genorm.Asc,
				},
			},
		},
		{
			description: "desc",
			addItem: orderItem{
				expr: &expr{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
				direction: genorm.Desc,
			},
			afterExprs: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
					direction: genorm.Desc,
				},
			},
		},
		{
			description: "empty before expr",
			beforeExprs: []orderItem{},
			addItem: orderItem{
				expr: &expr{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
				direction: genorm.Asc,
			},
			afterExprs: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
					direction: genorm.Asc,
				},
			},
		},
		{
			description: "condition already exists",
			beforeExprs: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
					direction: genorm.Asc,
				},
			},
			addItem: orderItem{
				expr: &expr{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(2)},
				},
				direction: genorm.Asc,
			},
			afterExprs: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
					direction: genorm.Asc,
				},
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(2)},
					},
					direction: genorm.Asc,
				},
			},
		},
		{
			description: "nil expr",
			addItem: orderItem{
				direction: genorm.Asc,
			},
			err: true,
		},
		{
			description: "invalid direction",
			addItem: orderItem{
				expr: &expr{
					query: "(`hoge`.`huga` = ?)",
					args:  []genorm.ExprType{genorm.Wrap(1)},
				},
			},
			err: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var beforeExprs []genorm.OrderItem[*mock.MockTable]
			if test.beforeExprs == nil {
				beforeExprs = nil
			} else {
				beforeExprs = make([]genorm.OrderItem[*mock.MockTable], 0, len(test.beforeExprs))
				for _, item := range test.beforeExprs {
					var expr genorm.TableExpr[*mock.MockTable]
					if item.expr == nil {
						expr = nil
					} else {
						mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
						expr = mockExpr
					}

					beforeExprs = append(beforeExprs, genorm.NewOrderItem(expr, item.direction))
				}
			}

			var afterExprs []genorm.OrderItem[*mock.MockTable]
			if test.afterExprs == nil {
				afterExprs = nil
			} else {
				afterExprs = make([]genorm.OrderItem[*mock.MockTable], 0, len(test.afterExprs))
				for _, item := range test.afterExprs {
					var expr genorm.TableExpr[*mock.MockTable]
					if item.expr == nil {
						expr = nil
					} else {
						mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
						expr = mockExpr
					}

					afterExprs = append(afterExprs, genorm.NewOrderItem(expr, item.direction))
				}
			}

			var expr genorm.TableExpr[*mock.MockTable]
			if test.addItem.expr == nil {
				expr = nil
			} else {
				mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
				expr = mockExpr
			}
			addItem := genorm.NewOrderItem(expr, test.addItem.direction)

			c := genorm.NewOrderClause(beforeExprs)

			err := c.Add(addItem)

			if test.err {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			afterItems := c.GetItems()
			assert.Equal(t, len(test.afterExprs), len(afterItems))

			for i, item := range afterItems {
				expr, direction := item.Value()
				expectExpr, expectDirection := afterExprs[i].Value()
				assert.Equal(t, expectExpr, expr)
				assert.Equal(t, expectDirection, direction)
			}
		})
	}
}

func TestOrderClauseExistTest(t *testing.T) {
	t.Parallel()

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	type orderItem struct {
		expr      *expr
		direction genorm.OrderDirection
	}

	tests := []struct {
		description string
		items       []orderItem
		exist       bool
	}{
		{
			description: "normal",
			items: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
					direction: genorm.Asc,
				},
			},
			exist: true,
		},
		{
			description: "nil expr",
		},
		{
			description: "no expr",
			items:       []orderItem{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var items []genorm.OrderItem[*mock.MockTable]
			if test.items == nil {
				items = nil
			} else {
				items = make([]genorm.OrderItem[*mock.MockTable], 0, len(test.items))
				for _, item := range test.items {
					var expr genorm.TableExpr[*mock.MockTable]
					if item.expr == nil {
						expr = nil
					} else {
						mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
						expr = mockExpr
					}

					items = append(items, genorm.NewOrderItem(expr, item.direction))
				}
			}

			c := genorm.NewOrderClause(items)

			res := c.Exists()

			assert.Equal(t, test.exist, res)
		})
	}
}

func TestOrderClauseGetExprTest(t *testing.T) {
	t.Parallel()

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	type orderItem struct {
		expr      *expr
		direction genorm.OrderDirection
	}

	tests := []struct {
		description string
		items       []orderItem
		query       string
		args        []genorm.ExprType
		err         bool
	}{
		{
			description: "normal",
			items: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
					direction: genorm.Asc,
				},
			},
			query: "ORDER BY (`hoge`.`huga` = ?) ASC",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "desc",
			items: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
					direction: genorm.Desc,
				},
			},
			query: "ORDER BY (`hoge`.`huga` = ?) DESC",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description: "invlid direction",
			items: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
				},
			},
			err: true,
		},
		{
			description: "multi exprs",
			items: []orderItem{
				{
					expr: &expr{
						query: "(`hoge`.`huga` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(1)},
					},
					direction: genorm.Asc,
				},
				{
					expr: &expr{
						query: "(`hoge`.`nya` = ?)",
						args:  []genorm.ExprType{genorm.Wrap(2)},
					},
					direction: genorm.Asc,
				},
			},
			query: "ORDER BY (`hoge`.`huga` = ?) ASC, (`hoge`.`nya` = ?) ASC",
			args:  []genorm.ExprType{genorm.Wrap(1), genorm.Wrap(2)},
		},
		{
			description: "nil expr",
			err:         true,
		},
		{
			description: "no expr",
			items:       []orderItem{},
			err:         true,
		},
		{
			description: "error condition",
			items: []orderItem{
				{
					expr: &expr{
						errs: []error{errors.New("error")},
					},
					direction: genorm.Asc,
				},
			},
			err: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var items []genorm.OrderItem[*mock.MockTable]
			if test.items == nil {
				items = nil
			} else {
				items = make([]genorm.OrderItem[*mock.MockTable], 0, len(test.items))
				for _, item := range test.items {
					var expr genorm.TableExpr[*mock.MockTable]
					if item.expr == nil {
						expr = nil
					} else {
						mockExpr := mock.NewMockTypedTableExpr[*mock.MockTable, genorm.WrappedPrimitive[bool]](ctrl)
						expr = mockExpr

						mockExpr.
							EXPECT().
							Expr().
							Return(item.expr.query, item.expr.args, item.expr.errs)
					}

					items = append(items, genorm.NewOrderItem(expr, item.direction))
				}
			}

			c := genorm.NewOrderClause(items)

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

func TestLimitClauseSetTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		before      uint64
		set         uint64
		err         bool
	}{
		{
			description: "normal",
			set:         1,
		},
		{
			description: "limit already set",
			before:      1,
			set:         2,
			err:         true,
		},
		{
			description: "limit 0",
			set:         0,
			err:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := genorm.NewLimitClause(test.before)

			err := c.Set(test.set)

			if test.err {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			assert.Equal(t, test.set, c.GetLimit())
		})
	}
}

func TestLimitClauseExistTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		limit       uint64
		exist       bool
	}{
		{
			description: "normal",
			limit:       1,
			exist:       true,
		},
		{
			description: "not exist",
			limit:       0,
			exist:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := genorm.NewLimitClause(test.limit)

			res := c.Exists()

			assert.Equal(t, test.exist, res)
		})
	}
}

func TestLimitClauseGetExprTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		limit       uint64
		query       string
		args        []genorm.ExprType
		err         bool
	}{
		{
			description: "normal",
			limit:       1,
			query:       "LIMIT 1",
		},
		{
			description: "empty limit",
			err:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := genorm.NewLimitClause(test.limit)

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

func TestOffsetClauseSetTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		before      uint64
		set         uint64
		err         bool
	}{
		{
			description: "normal",
			set:         1,
		},
		{
			description: "offset already set",
			before:      1,
			set:         2,
			err:         true,
		},
		{
			description: "limit 0",
			set:         0,
			err:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := genorm.NewOffsetClause(test.before)

			err := c.Set(test.set)

			if test.err {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			assert.Equal(t, test.set, c.GetOffset())
		})
	}
}

func TestOffsetClauseExistTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		offset      uint64
		exist       bool
	}{
		{
			description: "normal",
			offset:      1,
			exist:       true,
		},
		{
			description: "not exist",
			offset:      0,
			exist:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := genorm.NewOffsetClause(test.offset)

			res := c.Exists()

			assert.Equal(t, test.exist, res)
		})
	}
}

func TestOffsetClauseGetExprTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		offset      uint64
		query       string
		args        []genorm.ExprType
		err         bool
	}{
		{
			description: "normal",
			offset:      1,
			query:       "OFFSET 1",
		},
		{
			description: "empty limit",
			err:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := genorm.NewOffsetClause(test.offset)

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

func TestLockClauseSetTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		before      genorm.LockType
		set         genorm.LockType
		err         bool
	}{
		{
			description: "normal",
			set:         genorm.ForUpdate,
		},
		{
			description: "for share",
			set:         genorm.ForShare,
		},
		{
			description: "invalid lock type",
			set:         100,
			err:         true,
		},
		{
			description: "lockType already set",
			before:      genorm.ForUpdate,
			set:         genorm.ForUpdate,
			err:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := genorm.NewLockClause(test.before)

			err := c.Set(test.set)

			if test.err {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			assert.Equal(t, test.set, c.GetLockType())
		})
	}
}

func TestLockClauseExistTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		lockType    genorm.LockType
		exist       bool
	}{
		{
			description: "normal",
			lockType:    genorm.ForUpdate,
			exist:       true,
		},
		{
			description: "for share",
			lockType:    genorm.ForShare,
			exist:       true,
		},
		{
			description: "not exist",
			lockType:    0,
			exist:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := genorm.NewLockClause(test.lockType)

			res := c.Exists()

			assert.Equal(t, test.exist, res)
		})
	}
}

func TestLockClauseGetExprTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		lockType    genorm.LockType
		query       string
		args        []genorm.ExprType
		err         bool
	}{
		{
			description: "normal",
			lockType:    genorm.ForUpdate,
			query:       "FOR UPDATE",
		},
		{
			description: "for share",
			lockType:    genorm.ForShare,
			query:       "FOR SHARE",
		},
		{
			description: "invalid lock type",
			lockType:    100,
			err:         true,
		},
		{
			description: "empty lock type",
			err:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := genorm.NewLockClause(test.lockType)

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
