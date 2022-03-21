package relation

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/mock"
	"github.com/stretchr/testify/assert"
)

func TestJoinedTableName(t *testing.T) {
	t.Parallel()

	type expr struct {
		query string
		args  []genorm.ExprType
		errs  []error
	}

	tests := []struct {
		description  string
		relationType RelationType
		baseExpr     expr
		refTableExpr expr
		onExpr       *expr
		query        string
		args         []genorm.ExprType
		err          bool
	}{
		{
			description:  "normal",
			relationType: join,
			baseExpr: expr{
				query: "hoge",
			},
			refTableExpr: expr{
				query: "fuga",
			},
			query: "(hoge CROSS JOIN fuga)",
			args:  []genorm.ExprType{},
		},
		{
			description:  "onExpr",
			relationType: join,
			baseExpr: expr{
				query: "hoge",
			},
			refTableExpr: expr{
				query: "fuga",
			},
			onExpr: &expr{
				query: "(hoge.id = fuga.id)",
			},
			query: "(hoge INNER JOIN fuga ON (hoge.id = fuga.id))",
			args:  []genorm.ExprType{},
		},
		{
			description:  "onExpr with args",
			relationType: join,
			baseExpr: expr{
				query: "hoge",
			},
			refTableExpr: expr{
				query: "fuga",
			},
			onExpr: &expr{
				query: "(hoge.id = ?)",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			query: "(hoge INNER JOIN fuga ON (hoge.id = ?))",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description:  "onExpr with error",
			relationType: join,
			baseExpr: expr{
				query: "hoge",
			},
			refTableExpr: expr{
				query: "fuga",
			},
			onExpr: &expr{
				errs: []error{errors.New("onExpr error")},
			},
			err: true,
		},
		{
			description:  "baseTable with args",
			relationType: join,
			baseExpr: expr{
				query: "(hoge INNER JOIN fuga ON (fuga.id = ?))",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			refTableExpr: expr{
				query: "fuga",
			},
			query: "((hoge INNER JOIN fuga ON (fuga.id = ?)) CROSS JOIN fuga)",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description:  "baseTable with error",
			relationType: join,
			baseExpr: expr{
				errs: []error{errors.New("error")},
			},
			refTableExpr: expr{
				query: "fuga",
			},
			err: true,
		},
		{
			description:  "refTable with args",
			relationType: join,
			baseExpr: expr{
				query: "hoge",
			},
			refTableExpr: expr{
				query: "(fuga INNER JOIN piyo ON (piyo.id = ?))",
				args:  []genorm.ExprType{genorm.Wrap(1)},
			},
			query: "(hoge CROSS JOIN (fuga INNER JOIN piyo ON (piyo.id = ?)))",
			args:  []genorm.ExprType{genorm.Wrap(1)},
		},
		{
			description:  "refTable with error",
			relationType: join,
			baseExpr: expr{
				query: "hoge",
			},
			refTableExpr: expr{
				errs: []error{errors.New("error")},
			},
			err: true,
		},
		{
			description:  "left join",
			relationType: leftJoin,
			baseExpr: expr{
				query: "hoge",
			},
			refTableExpr: expr{
				query: "fuga",
			},
			onExpr: &expr{
				query: "(hoge.id = fuga.id)",
			},
			query: "(hoge LEFT JOIN fuga ON (hoge.id = fuga.id))",
			args:  []genorm.ExprType{},
		},
		{
			description:  "right join",
			relationType: rightJoin,
			baseExpr: expr{
				query: "hoge",
			},
			refTableExpr: expr{
				query: "fuga",
			},
			onExpr: &expr{
				query: "(hoge.id = fuga.id)",
			},
			query: "(hoge RIGHT JOIN fuga ON (hoge.id = fuga.id))",
			args:  []genorm.ExprType{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			baseTable := mock.NewMockTable(ctrl)
			baseTable.
				EXPECT().
				Expr().
				Return(test.baseExpr.query, test.baseExpr.args, test.baseExpr.errs)

			refTable := mock.NewMockTable(ctrl)
			if len(test.baseExpr.errs) == 0 {
				refTable.
					EXPECT().
					Expr().
					Return(test.refTableExpr.query, test.refTableExpr.args, test.refTableExpr.errs)
			}

			var onExpr genorm.Expr
			if test.onExpr != nil {
				mockExpr := mock.NewMockExpr(ctrl)
				mockExpr.
					EXPECT().
					Expr().
					Return(test.onExpr.query, test.onExpr.args, test.onExpr.errs)

				onExpr = mockExpr
			}

			relation := &Relation{
				relationType: test.relationType,
				baseTable:    baseTable,
				refTable:     refTable,
				onExpr:       onExpr,
			}

			query, args, errs := relation.JoinedTableName()

			if test.err {
				assert.Greater(t, len(errs), 0)
				return
			} else {
				if !assert.Len(t, errs, 0) {
					return
				}
			}

			assert.Equal(t, test.query, query)
			assert.Equal(t, test.args, args)
		})
	}
}
