package codegen

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mazrean/genorm/cmd/generator/types"
)

type column struct {
	table                 *table
	columnName            string
	fieldIdent            *ast.Ident
	fieldType             ast.Expr
	typeIdent             *ast.Ident
	varIdent              *ast.Ident
	tablePackageVarIdent  *ast.Ident
	tablePackageExprIdent *ast.Ident
	recvIdent             *ast.Ident
}

func newColumn(tbl *table, clmn *types.Column) *column {
	return &column{
		table:                 tbl,
		columnName:            clmn.Name,
		fieldIdent:            ast.NewIdent(clmn.FieldName),
		fieldType:             fieldTypeExpr(clmn.Type),
		typeIdent:             ast.NewIdent(tbl.lowerName() + clmn.FieldName),
		varIdent:              ast.NewIdent(tbl.name + clmn.FieldName),
		tablePackageVarIdent:  ast.NewIdent(clmn.FieldName),
		tablePackageExprIdent: ast.NewIdent(clmn.FieldName + "Expr"),
		recvIdent:             ast.NewIdent("c"),
	}
}

func (clmn *column) field() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{clmn.fieldIdent},
		Type:  clmn.fieldType,
	}
}

func fieldTypeExpr(columnTypeExpr ast.Expr) ast.Expr {
	columnTypeIdentExpr, ok := columnTypeExpr.(*ast.Ident)
	if ok {
		switch columnTypeIdentExpr.Name {
		case "bool",
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64",
			"string":
			return wrappedPrimitive(columnTypeExpr)
		default:
			return columnTypeExpr
		}
	}

	columnTypeSelectorExpr, ok := columnTypeExpr.(*ast.SelectorExpr)
	if ok && columnTypeSelectorExpr != nil {
		identExpr, ok := columnTypeSelectorExpr.X.(*ast.Ident)
		if ok &&
			identExpr != nil &&
			identExpr.Name == "time" &&
			columnTypeSelectorExpr.Sel.Name != "Time" {
			return wrappedPrimitive(columnTypeExpr)
		}
	}

	return columnTypeExpr
}

func (clmn *column) decls() []ast.Decl {
	return []ast.Decl{
		clmn.structDecl(),
		clmn.varDecl(),
		clmn.exprDecl(),
		clmn.sqlColumnsDecl(),
		clmn.tableNameDecl(),
		clmn.columnNameDecl(),
		clmn.tableExprDecl(),
		clmn.typeExprDecl(),
	}
}

func (clmn *column) structDecl() ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: clmn.typeIdent,
				Type: &ast.StructType{
					Fields: &ast.FieldList{},
				},
			},
		},
	}
}

func (clmn *column) varDecl() ast.Decl {
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{clmn.varIdent},
				Type: typedTableColumn(&ast.StarExpr{
					X: clmn.table.structIdent,
				}, clmn.fieldType),
				Values: []ast.Expr{
					&ast.CompositeLit{
						Type: clmn.typeIdent,
					},
				},
			},
		},
	}
}

func (clmn *column) exprDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{clmn.recvIdent},
					Type:  clmn.typeIdent,
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
								X:   clmn.recvIdent,
								Sel: columnSQLColumnsIdent,
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

func (clmn *column) sqlColumnsDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{clmn.recvIdent},
					Type:  clmn.typeIdent,
				},
			},
		},
		Name: columnSQLColumnsIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
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
									Value: "\"`%s`.`%s`\"",
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   clmn.recvIdent,
										Sel: columnTableNameIdent,
									},
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   clmn.recvIdent,
										Sel: columnColumnNameIdent,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (clmn *column) tableNameDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{clmn.recvIdent},
					Type:  clmn.typeIdent,
				},
			},
		},
		Name: columnTableNameIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
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
								X: &ast.CallExpr{
									Fun: &ast.ParenExpr{
										X: &ast.StarExpr{
											X: clmn.table.structIdent,
										},
									},
									Args: []ast.Expr{
										ast.NewIdent("nil"),
									},
								},
								Sel: columnTableNameIdent,
							},
						},
					},
				},
			},
		},
	}
}

func (clmn *column) columnNameDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{clmn.recvIdent},
					Type:  clmn.typeIdent,
				},
			},
		},
		Name: columnColumnNameIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s"`, clmn.columnName),
						},
					},
				},
			},
		},
	}
}

func (clmn *column) tableExprDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{clmn.recvIdent},
					Type:  clmn.typeIdent,
				},
			},
		},
		Name: tableExprTableExprIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: clmn.table.structIdent,
						},
					},
				},
			},
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
								X:   clmn.recvIdent,
								Sel: exprExprIdent,
							},
						},
					},
				},
			},
		},
	}
}

func (clmn *column) typeExprDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{clmn.recvIdent},
					Type:  clmn.typeIdent,
				},
			},
		},
		Name: typedExprTypedExprIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: clmn.fieldType,
					},
				},
			},
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
								X:   clmn.recvIdent,
								Sel: exprExprIdent,
							},
						},
					},
				},
			},
		},
	}
}

func (clmn *column) tablePackageDecls() []ast.Decl {
	return []ast.Decl{
		clmn.tablePackageVarDecl(),
		clmn.tablePackageExprDecl(),
	}
}

func (clmn *column) tablePackageVarDecl() ast.Decl {
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{clmn.tablePackageVarIdent},
				Type: typedTableColumn(&ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   rootPackageIdent,
						Sel: clmn.table.structIdent,
					},
				}, clmn.fieldType),
				Values: []ast.Expr{
					&ast.SelectorExpr{
						X:   rootPackageIdent,
						Sel: clmn.varIdent,
					},
				},
			},
		},
	}
}

func (clmn *column) tablePackageExprDecl() ast.Decl {
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{clmn.tablePackageExprIdent},
				Type: typedTableExpr(&ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   rootPackageIdent,
						Sel: clmn.table.structIdent,
					},
				}, clmn.fieldType),
				Values: []ast.Expr{
					&ast.SelectorExpr{
						X:   rootPackageIdent,
						Sel: clmn.varIdent,
					},
				},
			},
		},
	}
}
