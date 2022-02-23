package generator

import (
	"go/ast"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertTables(t *testing.T) {
	t.Parallel()

	typeIdent1 := ast.NewIdent("int64")
	typeIdent2 := ast.NewIdent("int")
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent("GetID"),
	}

	messageTable := &Table{
		StructName: "Message",
		Columns: []*Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent2,
			},
		},
		Methods:         []*Method{},
		RefTables:       []*RefTable{},
		RefJoinedTables: nil,
	}

	tests := []struct {
		description  string
		parserTables []*ParserTable
		tables       []*Table
		err          bool
	}{
		{
			description: "simple",
			parserTables: []*ParserTable{
				{
					StructName: "User",
					Columns: []*ParserColumn{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*ParserMethod{
						{
							StructName: "User",
							Type:       methodTypeIdentifier,
							Decl:       funcDecl,
						},
					},
					RefTables: []*ParserRefTable{},
				},
			},
			tables: []*Table{
				{
					StructName: "User",
					Columns: []*Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*Method{
						{
							Type: methodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables:       []*RefTable{},
					RefJoinedTables: nil,
				},
			},
		},
		{
			description: "reference",
			parserTables: []*ParserTable{
				{
					StructName: "User",
					Columns: []*ParserColumn{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*ParserMethod{
						{
							StructName: "User",
							Type:       methodTypeIdentifier,
							Decl:       funcDecl,
						},
					},
					RefTables: []*ParserRefTable{
						{
							FieldName:  "Message",
							StructName: "Message",
						},
					},
				},
				{
					StructName: "Message",
					Columns: []*ParserColumn{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent2,
						},
					},
					Methods:   []*ParserMethod{},
					RefTables: []*ParserRefTable{},
				},
			},
			tables: []*Table{
				{
					StructName: "User",
					Columns: []*Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*Method{
						{
							Type: methodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables: []*RefTable{
						{
							Table: messageTable,
						},
					},
					RefJoinedTables: nil,
				},
				messageTable,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			tables, err := convertTables(test.parserTables)
			if err != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.Len(t, tables, len(test.tables))

			for i, table := range tables {
				assert.Equal(t, test.tables[i], table)
			}
		})
	}
}

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

	messageTable := &Table{
		StructName: "Message",
		Columns: []*Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods:         []*Method{},
		RefTables:       []*RefTable{},
		RefJoinedTables: []*RefJoinedTable{},
	}
	userTable := &Table{
		StructName: "User",
		Columns: []*Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods: []*Method{
			{
				Type: methodTypeIdentifier,
				Decl: funcDecl,
			},
		},
		RefTables: []*RefTable{
			{
				Table: messageTable,
			},
		},
		RefJoinedTables: []*RefJoinedTable{},
	}

	messageOptionTable2 := &Table{
		StructName: "MessageOption",
		Columns: []*Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods:         []*Method{},
		RefTables:       []*RefTable{},
		RefJoinedTables: []*RefJoinedTable{},
	}
	messageTable2 := &Table{
		StructName: "Message",
		Columns: []*Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods: []*Method{},
		RefTables: []*RefTable{
			{
				Table: messageOptionTable2,
			},
		},
		RefJoinedTables: []*RefJoinedTable{},
	}
	userTable2 := &Table{
		StructName: "User",
		Columns: []*Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent1,
			},
		},
		Methods: []*Method{
			{
				Type: methodTypeIdentifier,
				Decl: funcDecl,
			},
		},
		RefTables: []*RefTable{
			{
				Table: messageTable2,
			},
		},
		RefJoinedTables: []*RefJoinedTable{
			{
				Table: &JoinedTable{
					Tables:          []*Table{messageTable2, messageOptionTable2},
					RefTables:       []*RefTable{},
					RefJoinedTables: []*RefJoinedTable{},
				},
			},
		},
	}

	tests := []struct {
		description        string
		tables             []*Table
		expectTables       []*Table
		expectJoinedTables []*JoinedTable
		joinNum            int
		err                bool
	}{
		{
			description: "simple",
			joinNum:     5,
			tables: []*Table{
				{
					StructName: "User",
					Columns: []*Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*Method{
						{
							Type: methodTypeIdentifier,
							Decl: funcDecl,
						},
					},
				},
			},
			expectTables: []*Table{
				{
					StructName: "User",
					Columns: []*Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*Method{
						{
							Type: methodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables:       []*RefTable{},
					RefJoinedTables: []*RefJoinedTable{},
				},
			},
			expectJoinedTables: []*JoinedTable{},
		},
		{
			description: "join",
			joinNum:     5,
			tables: []*Table{
				{
					StructName: "User",
					Columns: []*Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*Method{
						{
							Type: methodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables: []*RefTable{
						{
							Table: messageTable,
						},
					},
				},
				messageTable,
			},
			expectTables: []*Table{
				userTable,
				messageTable,
			},
			expectJoinedTables: []*JoinedTable{
				{
					Tables:          []*Table{messageTable, userTable},
					RefTables:       []*RefTable{},
					RefJoinedTables: []*RefJoinedTable{},
				},
			},
		},
		{
			description: "no join",
			joinNum:     1,
			tables: []*Table{
				{
					StructName: "User",
					Columns: []*Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*Method{
						{
							Type: methodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables: []*RefTable{
						{
							Table: messageTable,
						},
					},
				},
				messageTable,
			},
			expectTables: []*Table{
				userTable,
				messageTable,
			},
			expectJoinedTables: []*JoinedTable{},
		},
		{
			description: "multiple join",
			joinNum:     5,
			tables: []*Table{
				{
					StructName: "User",
					Columns: []*Column{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*Method{
						{
							Type: methodTypeIdentifier,
							Decl: funcDecl,
						},
					},
					RefTables: []*RefTable{
						{
							Table: messageTable2,
						},
					},
				},
				messageTable2,
				messageOptionTable2,
			},
			expectTables: []*Table{
				userTable2,
				messageTable2,
				messageOptionTable2,
			},
			expectJoinedTables: []*JoinedTable{
				{
					Tables: []*Table{messageTable2, userTable2},
					RefTables: []*RefTable{
						{
							Table: messageOptionTable2,
						},
					},
					RefJoinedTables: []*RefJoinedTable{},
				},
				{
					Tables:          []*Table{messageTable2, messageOptionTable2},
					RefTables:       []*RefTable{},
					RefJoinedTables: []*RefJoinedTable{},
				},
				{
					Tables:          []*Table{messageTable2, messageOptionTable2, userTable2},
					RefTables:       []*RefTable{},
					RefJoinedTables: []*RefJoinedTable{},
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
