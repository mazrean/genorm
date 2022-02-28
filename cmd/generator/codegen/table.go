package codegen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/mazrean/genorm/cmd/generator/types"
)

type table struct {
	table           *types.Table
	name            string
	structIdent     *ast.Ident
	recvIdent       *ast.Ident
	methods         []*method
	columns         []*column
	refTables       []*refTable
	refJoinedTables []*refJoinedTable
}

func newTable(tbl *types.Table) (*table, error) {
	codegenTable := &table{
		table:       tbl,
		name:        tbl.StructName,
		structIdent: ast.NewIdent(fmt.Sprintf("%sTable", tbl.StructName)),
		recvIdent:   ast.NewIdent("t"),
	}

	methods := make([]*method, 0, len(tbl.Methods))
	for _, m := range tbl.Methods {
		mthd, err := newMethod(codegenTable, m)
		if err != nil {
			return nil, fmt.Errorf("failed to create method: %w", err)
		}

		methods = append(methods, mthd)
	}
	codegenTable.methods = methods

	columns := make([]*column, 0, len(tbl.Columns))
	for _, c := range tbl.Columns {
		col := newColumn(codegenTable, c)

		columns = append(columns, col)
	}
	codegenTable.columns = columns

	return codegenTable, nil
}

func (tbl *table) lowerName() string {
	return strings.ToLower(tbl.name[0:1]) + tbl.name[1:]
}

func (tbl *table) decl() []ast.Decl {
	tableDecls := []ast.Decl{}

	tableDecls = append(tableDecls, tbl.structDecl())

	for _, method := range tbl.methods {
		tableDecls = append(tableDecls, method.Decl)
	}

	tableDecls = append(tableDecls, tbl.exprDecl(), tbl.columnsDecl(), tbl.columnMapDecl(), tbl.getErrorsDecl())

	for _, column := range tbl.columns {
		tableDecls = append(tableDecls, column.decls()...)
	}

	return tableDecls
}

func (tbl *table) structDecl() ast.Decl {
	fields := make([]*ast.Field, 0, len(tbl.columns))
	for _, column := range tbl.columns {
		fields = append(fields, column.field())
	}

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: tbl.structIdent,
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fields,
					},
				},
			},
		},
	}
}

func (tbl *table) exprDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{tbl.recvIdent},
					Type: &ast.StarExpr{
						X: tbl.structIdent,
					},
				},
			},
		},
		Name: exprExprIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
					},
					{
						Type: &ast.ArrayType{
							Elt: exprTypeInterfaceTypeExpr,
						},
					},
					{
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("error"),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   fmtIdent,
								Sel: ast.NewIdent("Sprintf"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"`%s`\"",
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   tbl.recvIdent,
										Sel: basicTableTableNameIdent,
									},
								},
							},
						},
						ast.NewIdent("nil"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (tbl *table) columnsDecl() ast.Decl {
	columnExprs := make([]ast.Expr, 0, len(tbl.columns))
	for _, column := range tbl.columns {
		columnExprs = append(columnExprs, column.varIdent)
	}

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{tbl.recvIdent},
					Type: &ast.StarExpr{
						X: tbl.structIdent,
					},
				},
			},
		},
		Name: tableColumnsIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: columnInterfaceTypeExpr,
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.ArrayType{
								Elt: columnInterfaceTypeExpr,
							},
							Elts: columnExprs,
						},
					},
				},
			},
		},
	}
}

func (tbl *table) columnMapDecl() ast.Decl {
	columnMapKeyValueExprs := make([]ast.Expr, 0, len(tbl.columns))
	for _, column := range tbl.columns {
		columnMapKeyValueExprs = append(columnMapKeyValueExprs, &ast.KeyValueExpr{
			Key: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   column.varIdent,
					Sel: ast.NewIdent("SQLColumnName"),
				},
			},
			Value: &ast.UnaryExpr{
				Op: token.AND,
				X: &ast.SelectorExpr{
					X:   tbl.recvIdent,
					Sel: column.fieldIdent,
				},
			},
		})
	}

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{tbl.recvIdent},
					Type: &ast.StarExpr{
						X: tbl.structIdent,
					},
				},
			},
		},
		Name: tableColumnMapIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.MapType{
							Key:   ast.NewIdent("string"),
							Value: columnFieldExprTypeExpr,
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.MapType{
								Key:   ast.NewIdent("string"),
								Value: columnFieldExprTypeExpr,
							},
							Elts: columnMapKeyValueExprs,
						},
					},
				},
			},
		},
	}
}

func (tbl *table) getErrorsDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{tbl.recvIdent},
					Type: &ast.StarExpr{
						X: tbl.structIdent,
					},
				},
			},
		},
		Name: tableGetErrorsIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("error"),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{ast.NewIdent("nil")},
				},
			},
		},
	}
}

type method struct {
	Type types.MethodType
	Decl *ast.FuncDecl
}

func newMethod(tbl *table, m *types.Method) (*method, error) {
	mthd := &method{
		Type: m.Type,
		Decl: m.Decl,
	}

	if mthd.Decl == nil ||
		mthd.Decl.Recv == nil ||
		len(mthd.Decl.Recv.List) == 0 ||
		mthd.Decl.Recv.List[0] == nil ||
		mthd.Decl.Recv.List[0].Type == nil {
		return nil, errors.New("invalid method")
	}
	switch mthd.Type {
	case types.MethodTypeIdentifier:
		mthd.Decl.Recv.List[0].Names = []*ast.Ident{
			ast.NewIdent(tbl.structIdent.Name),
		}
		mthd.Decl.Recv.List[0].Type = tbl.structIdent
	case types.MethodTypeStar:
		mthd.Decl.Recv.List[0].Names = []*ast.Ident{
			ast.NewIdent(tbl.structIdent.Name),
		}

		star, ok := mthd.Decl.Recv.List[0].Type.(*ast.StarExpr)
		if !ok || star == nil {
			return nil, errors.New("invalid method")
		}

		star.X = tbl.structIdent
	default:
		return nil, errors.New("unknown method type")
	}

	return mthd, nil
}
