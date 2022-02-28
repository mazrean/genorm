package codegen

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"github.com/mazrean/genorm/cmd/generator/types"
)

type joinedTable struct {
	joinedTable          *types.JoinedTable
	name                 string
	structIdent          *ast.Ident
	relationFieldIdent   *ast.Ident
	errsFieldIdent       *ast.Ident
	tablesInterfaceIdent *ast.Ident
	recvIdent            *ast.Ident
	columnTypeIdent      *ast.Ident
	columnTypeFieldIdent *ast.Ident
	columnTypeRecvIdent  *ast.Ident
	columnParseFuncIdent *ast.Ident
	tables               []*table
	refTables            []*refTable
	refJoinedTables      []*refJoinedTable
}

func newJoinedTable(jt *types.JoinedTable) *joinedTable {
	structName := joinedTableName(jt)
	structIdent := ast.NewIdent(structName)

	return &joinedTable{
		joinedTable:          jt,
		name:                 structName,
		structIdent:          structIdent,
		relationFieldIdent:   ast.NewIdent("relation"),
		errsFieldIdent:       ast.NewIdent("errs"),
		tablesInterfaceIdent: ast.NewIdent(structName + "Tables"),
		recvIdent:            ast.NewIdent("jt"),
		columnTypeIdent:      ast.NewIdent(structName + "ColumnType"),
		columnTypeFieldIdent: ast.NewIdent("columnType"),
		columnTypeRecvIdent:  ast.NewIdent("ct"),
		columnParseFuncIdent: ast.NewIdent(structName + "Parse"),
	}
}

func joinedTableName(jt *types.JoinedTable) string {
	sort.SliceStable(jt.Tables, func(i, j int) bool {
		return jt.Tables[i].StructName < jt.Tables[j].StructName
	})
	tableNames := make([]string, 0, len(jt.Tables))
	for _, table := range jt.Tables {
		tableNames = append(tableNames, table.StructName)
	}
	structName := strings.Join(tableNames, "") + "JoinedTable"

	return structName
}

func (jt *joinedTable) decl() []ast.Decl {
	return []ast.Decl{
		jt.structDecl(),
		jt.exprDecl(),
		jt.columnsDecl(),
		jt.columnMapDecl(),
		jt.baseTables(),
		jt.getErrorsDecl(),
		jt.addErrorDecl(),
		jt.tablesInterfaceDecl(),
		jt.columnParseFuncDecl(),
		jt.columnTypeDecl(),
		jt.columnTypeExprDecl(),
		jt.columnTypeSQLColumnDecl(),
		jt.columnTypeTableNameDecl(),
		jt.columnTypeColumnNameDecl(),
		jt.columnTypeTableExprDecl(),
		jt.columnTypeTypedExprDecl(),
	}
}

func (jt *joinedTable) structDecl() ast.Decl {
	fields := make([]*ast.Field, 0, len(jt.tables)+2)
	for _, table := range jt.tables {
		fields = append(fields, &ast.Field{
			Type: &ast.StarExpr{
				X: table.structIdent,
			},
		})
	}
	fields = append(fields, &ast.Field{
		Names: []*ast.Ident{jt.relationFieldIdent},
		Type: &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   genormRelationIdent,
				Sel: ast.NewIdent("Relation"),
			},
		},
	}, &ast.Field{
		Names: []*ast.Ident{jt.errsFieldIdent},
		Type: &ast.ArrayType{
			Elt: ast.NewIdent("error"),
		},
	})

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: jt.structIdent,
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fields,
					},
				},
			},
		},
	}
}

func (jt *joinedTable) exprDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.recvIdent},
					Type: &ast.StarExpr{
						X: jt.structIdent,
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
								X:   jt.recvIdent,
								Sel: exprExprIdent,
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) columnsDecl() ast.Decl {
	columnExprs := make([]ast.Expr, 0, len(jt.tables))
	for _, table := range jt.tables {
		for _, column := range table.columns {
			columnExprs = append(columnExprs, &ast.CallExpr{
				Fun:  jt.columnParseFuncIdent,
				Args: []ast.Expr{column.varIdent},
			})
		}
	}

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.recvIdent},
					Type: &ast.StarExpr{
						X: jt.structIdent,
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

func (jt *joinedTable) columnMapDecl() ast.Decl {
	columnMapsIdent := ast.NewIdent("columnMaps")
	columnMapExprs := make([]ast.Expr, 0, len(jt.tables))
	for _, table := range jt.tables {
		columnMapExprs = append(columnMapExprs, &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X:   jt.recvIdent,
					Sel: table.structIdent,
				},
				Sel: tableColumnMapIdent,
			},
		})
	}

	newColumnMapIdent := ast.NewIdent("newColumnMap")
	columnMapIdent := ast.NewIdent("columnMap")
	keyIdent := ast.NewIdent("k")
	exprIdent := ast.NewIdent("expr")

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.recvIdent},
					Type: &ast.StarExpr{
						X: jt.structIdent,
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
				&ast.AssignStmt{
					Lhs: []ast.Expr{columnMapsIdent},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.ArrayType{
								Elt: &ast.MapType{
									Key:   ast.NewIdent("string"),
									Value: columnFieldExprTypeExpr,
								},
							},
							Elts: columnMapExprs,
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{newColumnMapIdent},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.MapType{
								Key:   ast.NewIdent("string"),
								Value: columnFieldExprTypeExpr,
							},
						},
					},
				},
				&ast.RangeStmt{
					Key:   ast.NewIdent("_"),
					Value: columnMapIdent,
					Tok:   token.DEFINE,
					X:     columnMapsIdent,
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.RangeStmt{
								Key:   keyIdent,
								Value: exprIdent,
								Tok:   token.DEFINE,
								X:     columnMapIdent,
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.AssignStmt{
											Lhs: []ast.Expr{
												&ast.IndexExpr{
													X:     newColumnMapIdent,
													Index: keyIdent,
												},
											},
											Tok: token.ASSIGN,
											Rhs: []ast.Expr{
												exprIdent,
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						newColumnMapIdent,
					},
				},
			},
		},
	}
}

func (jt *joinedTable) baseTables() ast.Decl {
	baseTables := make([]ast.Expr, 0, len(jt.tables))
	for _, table := range jt.tables {
		baseTables = append(baseTables, &ast.SelectorExpr{
			X:   jt.recvIdent,
			Sel: table.structIdent,
		})
	}

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.recvIdent},
					Type: &ast.StarExpr{
						X: jt.structIdent,
					},
				},
			},
		},
		Name: joinedTableBaseTablesIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: basicTableTypeExpr,
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
								Elt: basicTableTypeExpr,
							},
							Elts: baseTables,
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) getErrorsDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.recvIdent},
					Type: &ast.StarExpr{
						X: jt.structIdent,
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
					Results: []ast.Expr{
						&ast.SelectorExpr{
							X:   jt.recvIdent,
							Sel: jt.errsFieldIdent,
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) addErrorDecl() ast.Decl {
	errIdent := ast.NewIdent("err")

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.recvIdent},
					Type: &ast.StarExpr{
						X: jt.structIdent,
					},
				},
			},
		},
		Name: tableAddErrorIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{errIdent},
						Type:  ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   jt.recvIdent,
							Sel: jt.errsFieldIdent,
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("append"),
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   jt.recvIdent,
									Sel: jt.errsFieldIdent,
								},
								errIdent,
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) tablesInterfaceDecl() ast.Decl {
	tableStructIdents := []ast.Expr{}
	for _, table := range jt.tables {
		tableStructIdents = append(tableStructIdents, &ast.StarExpr{
			X: table.structIdent,
		})
	}

	if len(tableStructIdents) == 0 {
		return nil
	}

	interfaceFieldType := tableStructIdents[0]
	for _, ident := range tableStructIdents[1:] {
		interfaceFieldType = &ast.BinaryExpr{
			X:  interfaceFieldType,
			Op: token.OR,
			Y:  ident,
		}
	}

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: jt.tablesInterfaceIdent,
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: interfaceFieldType,
							},
							{
								Type: tableTypeExpr,
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) columnParseFuncDecl() ast.Decl {
	tableTypeParamIdent := ast.NewIdent("S")
	exprTypeParamIdent := ast.NewIdent("T")

	columnParamIdent := ast.NewIdent("column")

	return &ast.FuncDecl{
		Name: jt.columnParseFuncIdent,
		Type: &ast.FuncType{
			TypeParams: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{tableTypeParamIdent},
						Type:  jt.tablesInterfaceIdent,
					},
					{
						Names: []*ast.Ident{exprTypeParamIdent},
						Type:  exprTypeInterfaceTypeExpr,
					},
				},
			},
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{columnParamIdent},
						Type:  typedTableColumn(tableTypeParamIdent, exprTypeParamIdent),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: typedTableColumn(&ast.StarExpr{
							X: jt.structIdent,
						}, exprTypeParamIdent),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.IndexListExpr{
								X: jt.columnTypeIdent,
								Indices: []ast.Expr{
									tableTypeParamIdent,
									exprTypeParamIdent,
								},
							},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key:   jt.columnTypeFieldIdent,
									Value: columnParamIdent,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) columnTypeDecl() ast.Decl {
	tableTypeParamIdent := ast.NewIdent("S")
	exprTypeParamIdent := ast.NewIdent("T")

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: jt.columnTypeIdent,
				TypeParams: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{tableTypeParamIdent},
							Type:  jt.tablesInterfaceIdent,
						},
						{
							Names: []*ast.Ident{exprTypeParamIdent},
							Type:  exprTypeInterfaceTypeExpr,
						},
					},
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{jt.columnTypeFieldIdent},
								Type:  typedTableColumn(tableTypeParamIdent, exprTypeParamIdent),
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) columnTypeExprDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.columnTypeRecvIdent},
					Type: &ast.IndexListExpr{
						X:       jt.columnTypeIdent,
						Indices: []ast.Expr{ast.NewIdent("_"), ast.NewIdent("_")},
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
								X: &ast.SelectorExpr{
									X:   jt.columnTypeRecvIdent,
									Sel: jt.columnTypeFieldIdent,
								},
								Sel: exprExprIdent,
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) columnTypeSQLColumnDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.columnTypeRecvIdent},
					Type: &ast.IndexListExpr{
						X:       jt.columnTypeIdent,
						Indices: []ast.Expr{ast.NewIdent("_"), ast.NewIdent("_")},
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
								X: &ast.SelectorExpr{
									X:   jt.columnTypeRecvIdent,
									Sel: jt.columnTypeFieldIdent,
								},
								Sel: columnSQLColumnsIdent,
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) columnTypeTableNameDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.columnTypeRecvIdent},
					Type: &ast.IndexListExpr{
						X:       jt.columnTypeIdent,
						Indices: []ast.Expr{ast.NewIdent("_"), ast.NewIdent("_")},
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
								X: &ast.SelectorExpr{
									X:   jt.columnTypeRecvIdent,
									Sel: jt.columnTypeFieldIdent,
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

func (jt *joinedTable) columnTypeColumnNameDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.columnTypeRecvIdent},
					Type: &ast.IndexListExpr{
						X:       jt.columnTypeIdent,
						Indices: []ast.Expr{ast.NewIdent("_"), ast.NewIdent("_")},
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
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   jt.columnTypeRecvIdent,
									Sel: jt.columnTypeFieldIdent,
								},
								Sel: columnColumnNameIdent,
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) columnTypeTableExprDecl() ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.columnTypeRecvIdent},
					Type: &ast.IndexListExpr{
						X:       jt.columnTypeIdent,
						Indices: []ast.Expr{ast.NewIdent("_"), ast.NewIdent("_")},
					},
				},
			},
		},
		Name: tableExprTableExprIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: jt.structIdent,
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
								X:   jt.columnTypeRecvIdent,
								Sel: exprExprIdent,
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) columnTypeTypedExprDecl() ast.Decl {
	exprTypeParamIdent := ast.NewIdent("T")

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{jt.columnTypeRecvIdent},
					Type: &ast.IndexListExpr{
						X:       jt.columnTypeIdent,
						Indices: []ast.Expr{ast.NewIdent("_"), exprTypeParamIdent},
					},
				},
			},
		},
		Name: typedExprTypedExprIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: exprTypeParamIdent,
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
								X:   jt.columnTypeRecvIdent,
								Sel: exprExprIdent,
							},
						},
					},
				},
			},
		},
	}
}
