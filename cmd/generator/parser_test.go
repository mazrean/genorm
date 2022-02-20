package generator

import (
	"go/ast"
	"go/parser"
	"testing"
)

func TestParseFuncDecl(t *testing.T) {
	t.Parallel()

	methodFunc := &ast.FuncDecl{
		Name: ast.NewIdent("main"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: ast.NewIdent("a"),
				},
			},
		},
	}

	pointerExpr, err := parser.ParseExpr(`*a`)
	if err != nil {
		t.Fatalf("failed to parse expression: %s", err)
	}

	pointerMethodFunc := &ast.FuncDecl{
		Name: ast.NewIdent("main"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: pointerExpr,
				},
			},
		},
	}

	tests := []struct {
		description string
		f           *ast.FuncDecl
		method      *ParserMethod
		isMethod    bool
		err         bool
	}{
		{
			description: "normal func -> not method",
			f: &ast.FuncDecl{
				Name: ast.NewIdent("main"),
			},
		},
		{
			description: "method func -> method",
			f:           methodFunc,
			method: &ParserMethod{
				StructName: "a",
				Type:       methodTypeIdentifier,
				Decl:       methodFunc,
			},
			isMethod: true,
		},
		{
			description: "func(field list length is 0) -> not method",
			f: &ast.FuncDecl{
				Name: ast.NewIdent("main"),
				Recv: &ast.FieldList{
					List: []*ast.Field{},
				},
			},
		},
		{
			description: "func(field list length is nil) -> not method",
			f: &ast.FuncDecl{
				Name: ast.NewIdent("main"),
				Recv: &ast.FieldList{},
			},
		},
		{
			description: "pointer method func -> method",
			f:           pointerMethodFunc,
			method: &ParserMethod{
				StructName: "a",
				Type:       methodTypeStar,
				Decl:       pointerMethodFunc,
			},
			isMethod: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			method, isMethod, err := parseFuncDecl(test.f)
			if err != nil {
				if !test.err {
					t.Fatalf("failed to parse func decl: %s", err)
				}
				return
			}

			if test.err {
				t.Fatalf("expected error but got no error")
			}

			if isMethod != test.isMethod {
				t.Fatalf("is method is not match")
			}

			if method == nil {
				if test.method != nil {
					t.Fatalf("method is nil")
				}
				return
			}

			if method.Type != test.method.Type {
				t.Fatalf("method type is not match")
			}

			if method.StructName != test.method.StructName {
				t.Fatalf("struct name is not match")
			}

			if method.Decl != test.method.Decl {
				t.Fatalf("method decl is not match")
			}
		})
	}
}

func TestParseGenDecl(t *testing.T) {
	t.Parallel()

	fieldType := ast.NewIdent("string")

	tests := []struct {
		description string
		g           *ast.GenDecl
		tables      []*ParserTable
		err         bool
	}{
		{
			description: "normal gen decl -> success",
			g: &ast.GenDecl{
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("a"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("s"),
										},
										Type: fieldType,
									},
								},
							},
						},
					},
				},
			},
			tables: []*ParserTable{
				{
					StructName: "a",
					Columns: []*ParserColumn{
						{
							Name:      "s",
							FieldName: "s",
							Type:      fieldType,
						},
					},
					RefTables: []*ParserRefTable{},
				},
			},
		},
		{
			description: "skip non-type spec",
			g: &ast.GenDecl{
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("a"),
						},
						Values: []ast.Expr{
							ast.NewIdent("b"),
						},
						Type: fieldType,
					},
				},
			},
			tables: []*ParserTable{},
		},
		{
			description: "skip non-struct type spec",
			g: &ast.GenDecl{
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("a"),
						Type: &ast.ChanType{},
					},
				},
			},
			tables: []*ParserTable{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			tables, err := parseGenDecl(test.g)
			if err != nil {
				if !test.err {
					t.Fatalf("failed to parse gen decl: %s", err)
				}
				return
			}

			if test.err {
				t.Fatalf("expected error but got no error")
			}

			if len(tables) != len(test.tables) {
				t.Fatalf("table length is not match")
			}

			for i, table := range tables {
				if table.StructName != test.tables[i].StructName {
					t.Fatalf("struct name is not match")
				}

				if len(table.Columns) != len(test.tables[i].Columns) {
					t.Fatalf("column length is not match")
				}

				for j, column := range table.Columns {
					if column.Name != test.tables[i].Columns[j].Name {
						t.Fatalf("column name is not match")
					}

					if column.FieldName != test.tables[i].Columns[j].FieldName {
						t.Fatalf("column field name is not match")
					}

					if column.Type != test.tables[i].Columns[j].Type {
						t.Fatalf("column type is not match")
					}
				}

				if len(table.RefTables) != len(test.tables[i].RefTables) {
					t.Fatalf("ref table length is not match")
				}
			}
		})
	}
}
