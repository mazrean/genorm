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
