package convert

import (
	"go/ast"
	"sort"
	"testing"

	"github.com/mazrean/genorm/cmd/generator/types"
	"github.com/stretchr/testify/assert"
)

func TestTablesHash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		joinedTable *converterJoinedTable
		tableLength int
		hash        int64
	}{
		{
			description: "simple",
			joinedTable: &converterJoinedTable{
				hash: -1,
				tables: map[int]*converterTable{
					1: {},
				},
			},
			tableLength: 2,
			hash:        1,
		},
		{
			description: "id: 0",
			joinedTable: &converterJoinedTable{
				hash: -1,
				tables: map[int]*converterTable{
					0: {},
				},
			},
			tableLength: 2,
			hash:        0,
		},
		{
			description: "multiple",
			joinedTable: &converterJoinedTable{
				hash: -1,
				tables: map[int]*converterTable{
					1: {},
					3: {},
				},
			},
			tableLength: 4,
			hash:        13,
		},
		{
			description: "length: 0",
			joinedTable: &converterJoinedTable{
				hash:   -1,
				tables: map[int]*converterTable{},
			},
			tableLength: 0,
			hash:        0,
		},
		{
			description: "use cache",
			joinedTable: &converterJoinedTable{
				hash: 50,
			},
			tableLength: 10,
			hash:        50,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			assert.Equal(t, test.hash, test.joinedTable.tablesHash(test.tableLength))
		})
	}
}

func TestConvertJoinedTables(t *testing.T) {
	t.Parallel()

	typeIdent1 := ast.NewIdent("int64")
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent("GetID"),
	}

	messageTable := &types.Table{
		StructName: "Message",
		Columns: []*types.Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods:         []*types.Method{},
		RefTables:       []*types.RefTable{},
		RefJoinedTables: []*types.RefJoinedTable{},
	}
	userTable := &types.Table{
		StructName: "User",
		Columns: []*types.Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods: []*types.Method{
			{
				Type: types.MethodTypeIdentifier,
				Decl: funcDecl,
			},
		},
		RefTables: []*types.RefTable{
			{
				Table: messageTable,
			},
		},
		RefJoinedTables: []*types.RefJoinedTable{},
	}

	messageOptionTable2 := &types.Table{
		StructName: "MessageOption",
		Columns: []*types.Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods:         []*types.Method{},
		RefTables:       []*types.RefTable{},
		RefJoinedTables: []*types.RefJoinedTable{},
	}
	messageTable2 := &types.Table{
		StructName: "Message",
		Columns: []*types.Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods: []*types.Method{},
		RefTables: []*types.RefTable{
			{
				Table: messageOptionTable2,
			},
		},
		RefJoinedTables: []*types.RefJoinedTable{},
	}
	userTable2 := &types.Table{
		StructName: "User",
		Columns: []*types.Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods: []*types.Method{
			{
				Type: types.MethodTypeIdentifier,
				Decl: funcDecl,
			},
		},
		RefTables: []*types.RefTable{
			{
				Table: messageTable2,
			},
		},
		RefJoinedTables: []*types.RefJoinedTable{
			{
				Table: &types.JoinedTable{
					Tables:          []*types.Table{messageTable2, messageOptionTable2},
					RefTables:       []*types.RefTable{},
					RefJoinedTables: []*types.RefJoinedTable{},
				},
			},
		},
	}

	tests := []struct {
		description        string
		tables             []*types.Table
		expectTables       []*types.Table
		expectJoinedTables []*types.JoinedTable
		joinNum            int
		err                bool
	}{
		{
			description: "simple",
			joinNum:     5,
			tables: []*types.Table{
				{
					StructName: "User",
					Columns: []*types.Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*types.Method{
						{
							Type: types.MethodTypeIdentifier,
							Decl: funcDecl,
						},
					},
				},
			},
			expectTables: []*types.Table{
				{
					StructName: "User",
					Columns: []*types.Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*types.Method{
						{
							Type: types.MethodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables:       []*types.RefTable{},
					RefJoinedTables: []*types.RefJoinedTable{},
				},
			},
			expectJoinedTables: []*types.JoinedTable{},
		},
		{
			description: "join",
			joinNum:     5,
			tables: []*types.Table{
				{
					StructName: "User",
					Columns: []*types.Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*types.Method{
						{
							Type: types.MethodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables: []*types.RefTable{
						{
							Table: messageTable,
						},
					},
				},
				messageTable,
			},
			expectTables: []*types.Table{
				userTable,
				messageTable,
			},
			expectJoinedTables: []*types.JoinedTable{
				{
					Tables:          []*types.Table{messageTable, userTable},
					RefTables:       []*types.RefTable{},
					RefJoinedTables: []*types.RefJoinedTable{},
				},
			},
		},
		{
			description: "no join",
			joinNum:     1,
			tables: []*types.Table{
				{
					StructName: "User",
					Columns: []*types.Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*types.Method{
						{
							Type: types.MethodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables: []*types.RefTable{
						{
							Table: messageTable,
						},
					},
				},
				messageTable,
			},
			expectTables: []*types.Table{
				userTable,
				messageTable,
			},
			expectJoinedTables: []*types.JoinedTable{},
		},
		{
			description: "multiple join",
			joinNum:     5,
			tables: []*types.Table{
				{
					StructName: "User",
					Columns: []*types.Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*types.Method{
						{
							Type: types.MethodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables: []*types.RefTable{
						{
							Table: messageTable2,
						},
					},
				},
				messageTable2,
				messageOptionTable2,
			},
			expectTables: []*types.Table{
				userTable2,
				messageTable2,
				messageOptionTable2,
			},
			expectJoinedTables: []*types.JoinedTable{
				{
					Tables: []*types.Table{messageTable2, userTable2},
					RefTables: []*types.RefTable{
						{
							Table: messageOptionTable2,
						},
					},
					RefJoinedTables: []*types.RefJoinedTable{},
				},
				{
					Tables:          []*types.Table{messageTable2, messageOptionTable2},
					RefTables:       []*types.RefTable{},
					RefJoinedTables: []*types.RefJoinedTable{},
				},
				{
					Tables:          []*types.Table{messageTable2, messageOptionTable2, userTable2},
					RefTables:       []*types.RefTable{},
					RefJoinedTables: []*types.RefJoinedTable{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			tables, joinedTables, err := convertJoinedTables(test.tables, test.joinNum)
			if err != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.ElementsMatch(t, test.expectTables, tables)

			for _, joinedTable := range joinedTables {
				sort.Slice(joinedTable.Tables, func(i, j int) bool {
					return joinedTable.Tables[i].StructName < joinedTable.Tables[j].StructName
				})
			}
			assert.ElementsMatch(t, test.expectJoinedTables, joinedTables)
		})
	}
}
