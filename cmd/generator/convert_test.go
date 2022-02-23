package generator

import (
	"go/ast"
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
		joinedTable *generateJoinedTable
		tableLength int
		hash        int64
	}{
		{
			description: "simple",
			joinedTable: &generateJoinedTable{
				hash: -1,
				tables: map[int]*generateTable{
					1: {},
				},
			},
			tableLength: 2,
			hash:        1,
		},
		{
			description: "id: 0",
			joinedTable: &generateJoinedTable{
				hash: -1,
				tables: map[int]*generateTable{
					0: {},
				},
			},
			tableLength: 2,
			hash:        0,
		},
		{
			description: "multiple",
			joinedTable: &generateJoinedTable{
				hash: -1,
				tables: map[int]*generateTable{
					1: {},
					3: {},
				},
			},
			tableLength: 4,
			hash:        13,
		},
		{
			description: "length: 0",
			joinedTable: &generateJoinedTable{
				hash:   -1,
				tables: map[int]*generateTable{},
			},
			tableLength: 0,
			hash:        0,
		},
		{
			description: "use cache",
			joinedTable: &generateJoinedTable{
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
