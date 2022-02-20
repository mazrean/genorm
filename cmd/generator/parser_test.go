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
