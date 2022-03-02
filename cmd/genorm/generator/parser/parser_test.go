package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/mazrean/genorm/cmd/genorm/generator/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertTables(t *testing.T) {
	t.Parallel()

	typeIdent1 := ast.NewIdent("int64")
	typeIdent2 := ast.NewIdent("int")
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent("GetID"),
	}

	messageTable := &types.Table{
		StructName: "Message",
		Columns: []*types.Column{
			{
				Name:      "id",
				FieldName: "ID",
				Type:      typeIdent2,
			},
		},
		Methods:         []*types.Method{},
		RefTables:       []*types.RefTable{},
		RefJoinedTables: nil,
	}

	tests := []struct {
		description  string
		parserTables []*parserTable
		tables       []*types.Table
		err          bool
	}{
		{
			description: "simple",
			parserTables: []*parserTable{
				{
					StructName: "User",
					Columns: []*parserColumn{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*parserMethod{
						{
							StructName: "User",
							Type:       types.MethodTypeIdentifier,
							Decl:       funcDecl,
						},
					},
					RefTables: []*parserRefTable{},
				},
			},
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
					RefTables:       []*types.RefTable{},
					RefJoinedTables: nil,
				},
			},
		},
		{
			description: "reference",
			parserTables: []*parserTable{
				{
					StructName: "User",
					Columns: []*parserColumn{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent1,
						},
					},
					Methods: []*parserMethod{
						{
							StructName: "User",
							Type:       types.MethodTypeIdentifier,
							Decl:       funcDecl,
						},
					},
					RefTables: []*parserRefTable{
						{
							FieldName:  "Message",
							StructName: "Message",
						},
					},
				},
				{
					StructName: "Message",
					Columns: []*parserColumn{
						{
							Name:      "id",
							FieldName: "ID",
							Type:      typeIdent2,
						},
					},
					Methods:   []*parserMethod{},
					RefTables: []*parserRefTable{},
				},
			},
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

			assert.ElementsMatch(t, test.tables, tables)
		})
	}
}

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
		method      *parserMethod
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
			method: &parserMethod{
				StructName: "a",
				Type:       types.MethodTypeIdentifier,
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
			method: &parserMethod{
				StructName: "a",
				Type:       types.MethodTypeStar,
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
		tables      []*parserTable
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
			tables: []*parserTable{
				{
					StructName: "a",
					Columns: []*parserColumn{
						{
							Name:      "s",
							FieldName: "s",
							Type:      fieldType,
						},
					},
					RefTables: []*parserRefTable{},
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
			tables: []*parserTable{},
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
			tables: []*parserTable{},
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

func TestParseStructType(t *testing.T) {
	t.Parallel()

	fieldType := ast.NewIdent("string")
	refType := &ast.IndexExpr{
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("genorm"),
			Sel: ast.NewIdent("Ref"),
		},
		Index: &ast.Ident{
			Name: "Table",
		},
	}

	tests := []struct {
		description string
		name        string
		s           *ast.StructType
		table       *parserTable
		err         bool
	}{
		{
			description: "normal struct type -> success",
			name:        "a",
			s: &ast.StructType{
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
			table: &parserTable{
				StructName: "a",
				Columns: []*parserColumn{
					{
						Name:      "s",
						FieldName: "s",
						Type:      fieldType,
					},
				},
				RefTables: []*parserRefTable{},
			},
		},
		{
			description: "struct type(tag exist) -> success",
			name:        "a",
			s: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("s"),
							},
							Type: fieldType,
							Tag: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "`genorm:\"t\"`",
							},
						},
					},
				},
			},
			table: &parserTable{
				StructName: "a",
				Columns: []*parserColumn{
					{
						Name:      "t",
						FieldName: "s",
						Type:      fieldType,
					},
				},
				RefTables: []*parserRefTable{},
			},
		},
		{
			description: "struct type(ref exist) -> success",
			name:        "a",
			s: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("s"),
							},
							Type: refType,
						},
					},
				},
			},
			table: &parserTable{
				StructName: "a",
				Columns:    []*parserColumn{},
				RefTables: []*parserRefTable{
					{
						FieldName:  "s",
						StructName: "Table",
					},
				},
			},
		},
		{
			description: "struct type(multi column) -> success",
			name:        "a",
			s: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("s"),
							},
							Type: fieldType,
						},
						{
							Names: []*ast.Ident{
								ast.NewIdent("t"),
							},
							Type: fieldType,
						},
					},
				},
			},
			table: &parserTable{
				StructName: "a",
				Columns: []*parserColumn{
					{
						Name:      "s",
						FieldName: "s",
						Type:      fieldType,
					},
					{
						Name:      "t",
						FieldName: "t",
						Type:      fieldType,
					},
				},
				RefTables: []*parserRefTable{},
			},
		},
		{
			description: "struct type(multi name field) -> success",
			name:        "a",
			s: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("s"),
								ast.NewIdent("t"),
							},
							Type: fieldType,
						},
					},
				},
			},
			table: &parserTable{
				StructName: "a",
				Columns: []*parserColumn{
					{
						Name:      "s",
						FieldName: "s",
						Type:      fieldType,
					},
					{
						Name:      "t",
						FieldName: "t",
						Type:      fieldType,
					},
				},
				RefTables: []*parserRefTable{},
			},
		},
		{
			description: "struct type(multi name ref exist) -> success",
			name:        "a",
			s: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("s"),
								ast.NewIdent("t"),
							},
							Type: refType,
						},
					},
				},
			},
			table: &parserTable{
				StructName: "a",
				Columns:    []*parserColumn{},
				RefTables: []*parserRefTable{
					{
						FieldName:  "s",
						StructName: "Table",
					},
					{
						FieldName:  "t",
						StructName: "Table",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			table, err := parseStructType(test.name, test.s)
			if err != nil {
				if !test.err {
					t.Fatalf("failed to parse struct type: %v", err)
				}
				return
			}

			if test.err {
				t.Fatalf("expected error but got no error")
			}

			if table == nil {
				if test.table != nil {
					t.Fatalf("expected table but got nil")
				}
				return
			}

			if table.StructName != test.table.StructName {
				t.Fatalf("struct name is not match")
			}

			if len(table.Columns) != len(test.table.Columns) {
				t.Fatalf("column length is not match")
			}

			for j, column := range table.Columns {
				if column.Name != test.table.Columns[j].Name {
					t.Fatalf("column name is not match(expected: %s, actual: %s)", test.table.Columns[j].Name, column.Name)
				}

				if column.FieldName != test.table.Columns[j].FieldName {
					t.Fatalf("column field name is not match(expected: %s, actual: %s)", test.table.Columns[j].FieldName, column.FieldName)
				}

				if column.Type != test.table.Columns[j].Type {
					t.Fatalf("column type is not match(expected: %s, actual: %s)", test.table.Columns[j].Type, column.Type)
				}
			}

			if len(table.RefTables) != len(test.table.RefTables) {
				t.Fatalf("ref table length is not match")
			}

			for j, refTable := range table.RefTables {
				if refTable.FieldName != test.table.RefTables[j].FieldName {
					t.Fatalf("ref table field name is not match(expected: %s, actual: %s)", test.table.RefTables[j].FieldName, refTable.FieldName)
				}

				if refTable.StructName != test.table.RefTables[j].StructName {
					t.Fatalf("ref table struct name is not match(expected: %s, actual: %s)", test.table.RefTables[j].StructName, refTable.StructName)
				}
			}
		})
	}
}

func TestCheckRefType(t *testing.T) {
	fieldExpr1, err := parser.ParseExpr("string")
	if err != nil {
		t.Fatalf("failed to parse field expr: %v", err)
	}

	fieldExpr2, err := parser.ParseExpr("uuid.UUID")
	if err != nil {
		t.Fatalf("failed to parse field expr: %v", err)
	}

	fieldExpr3, err := parser.ParseExpr("genorm.Table")
	if err != nil {
		t.Fatalf("failed to parse field expr: %v", err)
	}

	fieldExpr4, err := parser.ParseExpr("uuid.UUID[Table]")
	if err != nil {
		t.Fatalf("failed to parse field expr: %v", err)
	}

	fieldExpr5, err := parser.ParseExpr("genorm.Table[Table]")
	if err != nil {
		t.Fatalf("failed to parse field expr: %v", err)
	}

	fieldExpr6, err := parser.ParseExpr("string[Table]")
	if err != nil {
		t.Fatalf("failed to parse field expr: %v", err)
	}

	refExpr, err := parser.ParseExpr("genorm.Ref[Table]")
	if err != nil {
		t.Fatalf("failed to parse field expr: %v", err)
	}

	tests := []struct {
		description string
		e           ast.Expr
		tableName   string
		isRef       bool
	}{
		{
			description: "string -> not ref",
			e:           fieldExpr1,
		},
		{
			description: "uuid.UUID -> not ref",
			e:           fieldExpr2,
		},
		{
			description: "genorm type -> not ref",
			e:           fieldExpr3,
		},
		{
			description: "uuid.UUID[Table] -> not ref",
			e:           fieldExpr4,
		},
		{
			description: "genorm type(with type param) -> not ref",
			e:           fieldExpr5,
		},
		{
			description: "string[Table] -> not ref",
			e:           fieldExpr6,
		},
		{
			description: "genorm.Ref[Table] -> ref",
			e:           refExpr,
			tableName:   "Table",
			isRef:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			tableName, isRef := checkRefType(test.e)

			if isRef != test.isRef {
				t.Fatalf("isRef is not match(expected: %t, actual: %t)", test.isRef, isRef)
			}

			if tableName != test.tableName {
				t.Fatalf("table name is not match(expected: %s, actual: %s)", test.tableName, tableName)
			}
		})
	}
}
