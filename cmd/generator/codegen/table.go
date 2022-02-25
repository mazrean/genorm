package codegen

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/mazrean/genorm/cmd/generator/types"
)

func tableDecl(table *types.Table) []ast.Decl {
	tableDecls := []ast.Decl{}

	structName := fmt.Sprintf("%sTable", table.StructName)
	structIdent := ast.NewIdent(structName)

	fields := make([]*ast.Field, 0, len(table.Columns))
	for _, column := range table.Columns {
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(column.FieldName),
			},
			Type: codegenFieldTypeExpr(column.Type),
		})
	}

	tableDecls = append(tableDecls, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: structIdent,
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fields,
					},
				},
			},
		},
	})

	for _, method := range table.Methods {
		method.SetStructName(structName)
		tableDecls = append(tableDecls, method.Decl)
	}

	columnExprs := make([]ast.Expr, 0, len(table.Columns))
	columnDecls := make([]ast.Decl, 0, len(table.Columns))
	for _, column := range table.Columns {
		columnExpr, newColumnDecls := columnMainDecl(table, column)
		columnExprs = append(columnExprs, columnExpr)
		columnDecls = append(columnDecls, newColumnDecls...)
	}

	recvIdent := ast.NewIdent("t")

	columnMapKeyValueExprs := make([]ast.Expr, 0, len(table.Columns))
	for i := range table.Columns {
		columnMapKeyValueExprs = append(columnMapKeyValueExprs, &ast.KeyValueExpr{
			Key: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   columnExprs[i],
					Sel: ast.NewIdent("SQLColumnName"),
				},
			},
			Value: &ast.UnaryExpr{
				Op: token.AND,
				X: &ast.SelectorExpr{
					X:   recvIdent,
					Sel: ast.NewIdent(table.Columns[i].FieldName),
				},
			},
		})
	}

	tableDecls = append(tableDecls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{recvIdent},
					Type: &ast.StarExpr{
						X: structIdent,
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
						Type: ast.NewIdent("error"),
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
										X:   recvIdent,
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
	}, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{recvIdent},
					Type: &ast.StarExpr{
						X: structIdent,
					},
				},
			},
		},
		Name: tableColumnsIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
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
	}, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{recvIdent},
					Type: &ast.StarExpr{
						X: structIdent,
					},
				},
			},
		},
		Name: tableColumnMapIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: &ast.MapType{
							Key:   ast.NewIdent("string"),
							Value: exprTypeInterfaceTypeExpr,
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
								Value: exprTypeInterfaceTypeExpr,
							},
							Elts: columnMapKeyValueExprs,
						},
					},
				},
			},
		},
	}, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{recvIdent},
					Type: &ast.StarExpr{
						X: structIdent,
					},
				},
			},
		},
		Name: tableGetErrors,
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
	})

	tableDecls = append(tableDecls, columnDecls...)

	return tableDecls
}

func columnMainDecl(table *types.Table, column *types.Column) (ast.Expr, []ast.Decl) {
	lowerTableName := strings.ToLower(table.StructName[0:1]) + table.StructName[1:]
	columnTypeIdent := ast.NewIdent(lowerTableName + column.FieldName)
	columnVarIdent := ast.NewIdent(table.StructName + column.FieldName)

	tableStructPointerType := &ast.UnaryExpr{
		Op: token.AND,
		X:  ast.NewIdent(fmt.Sprintf("%sTable", table.StructName)),
	}

	recvIdent := ast.NewIdent("c")
	columnDecls := make([]ast.Decl, 0, len(table.Columns)*8)
	for _, column := range table.Columns {
		columnDecls = append(columnDecls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: columnTypeIdent,
					Type: &ast.StructType{},
				},
			},
		}, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{recvIdent},
						Type: &ast.StarExpr{
							X: columnTypeIdent,
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
							Type: ast.NewIdent("error"),
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
									X:   recvIdent,
									Sel: columnSQLColumnsIdent,
								},
							},
							ast.NewIdent("nil"),
							ast.NewIdent("nil"),
						},
					},
				},
			},
		}, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{recvIdent},
						Type: &ast.StarExpr{
							X: columnTypeIdent,
						},
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
											X:   recvIdent,
											Sel: columnTableNameIdent,
										},
									},
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   recvIdent,
											Sel: columnColumnNameIdent,
										},
									},
								},
							},
						},
					},
				},
			},
		}, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{recvIdent},
						Type: &ast.StarExpr{
							X: columnTypeIdent,
						},
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
											X: tableStructPointerType,
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
		}, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{recvIdent},
						Type: &ast.StarExpr{
							X: columnTypeIdent,
						},
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
								Value: fmt.Sprintf(`"%s"`, column.Name),
							},
						},
					},
				},
			},
		}, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{recvIdent},
						Type: &ast.StarExpr{
							X: columnTypeIdent,
						},
					},
				},
			},
			Name: tableExprTableExprIdent,
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: tableStructPointerType,
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
							Type: ast.NewIdent("error"),
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
									X:   recvIdent,
									Sel: exprExprIdent,
								},
							},
						},
					},
				},
			},
		}, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{recvIdent},
						Type: &ast.StarExpr{
							X: columnTypeIdent,
						},
					},
				},
			},
			Name: typedExprTypedExprIdent,
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: codegenFieldTypeExpr(column.Type),
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
							Type: ast.NewIdent("error"),
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
									X:   recvIdent,
									Sel: exprExprIdent,
								},
							},
						},
					},
				},
			},
		}, &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{columnVarIdent},
					Type:  columnTypeIdent,
					Values: []ast.Expr{
						&ast.CompositeLit{
							Type: columnTypeIdent,
						},
					},
				},
			},
		}, &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{columnVarIdent},
					Type:  columnTypeIdent,
					Values: []ast.Expr{
						&ast.CompositeLit{
							Type: columnTypeIdent,
						},
					},
				},
			},
		})
	}

	return columnVarIdent, columnDecls
}
