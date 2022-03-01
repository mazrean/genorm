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
	name := joinedTableName(jt)
	structName := name + "JoinedTable"
	structIdent := ast.NewIdent(structName)

	return &joinedTable{
		joinedTable:          jt,
		name:                 name,
		structIdent:          structIdent,
		relationFieldIdent:   ast.NewIdent("relation"),
		errsFieldIdent:       ast.NewIdent("errs"),
		tablesInterfaceIdent: ast.NewIdent(name + "Tables"),
		recvIdent:            ast.NewIdent("jt"),
		columnTypeIdent:      ast.NewIdent(name + "ColumnType"),
		columnTypeFieldIdent: ast.NewIdent("columnType"),
		columnTypeRecvIdent:  ast.NewIdent("ct"),
		columnParseFuncIdent: ast.NewIdent(name + "Parse"),
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

	return strings.Join(tableNames, "")
}

func (jt *joinedTable) decl() []ast.Decl {
	decls := []ast.Decl{}

	decls = append(
		decls,
		jt.structDecl(),
		jt.exprDecl(),
		jt.columnsDecl(),
		jt.columnMapDecl(),
		jt.baseTables(),
		jt.getErrorsDecl(),
		jt.addErrorDecl(),
		jt.setRelationDecl(),
	)

	for _, ref := range jt.refTables {
		decls = append(decls, jt.tableJoinDecl(ref))
	}

	for _, ref := range jt.refJoinedTables {
		decls = append(decls, jt.joinedTableJoinDecl(ref))
	}

	decls = append(
		decls,
		jt.tablesInterfaceDecl(),
		jt.columnParseFuncDecl(),
		jt.columnTypeDecl(),
		jt.columnTypeExprDecl(),
		jt.columnTypeSQLColumnDecl(),
		jt.columnTypeTableNameDecl(),
		jt.columnTypeColumnNameDecl(),
		jt.columnTypeTableExprDecl(),
		jt.columnTypeTypedExprDecl(),
	)

	return decls
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
								X: &ast.SelectorExpr{
									X:   jt.recvIdent,
									Sel: jt.relationFieldIdent,
								},
								Sel: relationJoinedTableNameIdent,
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

func (jt *joinedTable) setRelationDecl() ast.Decl {
	relationIdent := ast.NewIdent("relation")

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
		Name: joinedTableSetRelationIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{relationIdent},
						Type: &ast.StarExpr{
							X: relationTypeExpr,
						},
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
							Sel: jt.relationFieldIdent,
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{relationIdent},
				},
			},
		},
	}
}

func (jt *joinedTable) tableJoinDecl(ref *refTable) ast.Decl {
	joinIdent := ast.NewIdent(ref.refTable.name)
	refIdent := ast.NewIdent("ref")

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
		Name: joinIdent,
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: relationContext(&ast.StarExpr{
							X: jt.structIdent,
						}, &ast.StarExpr{
							X: ref.refTable.structIdent,
						}, &ast.StarExpr{
							X: ref.joinedTable.structIdent,
						}),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{refIdent},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: ref.refTable.structIdent,
							Elts: []ast.Expr{},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: newRelationContext(&ast.StarExpr{
								X: jt.structIdent,
							}, &ast.StarExpr{
								X: ref.refTable.structIdent,
							}, &ast.StarExpr{
								X: ref.joinedTable.structIdent,
							}),
							Args: []ast.Expr{
								jt.recvIdent,
								&ast.UnaryExpr{
									Op: token.AND,
									X:  refIdent,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (jt *joinedTable) joinedTableJoinDecl(ref *refJoinedTable) ast.Decl {
	joinIdent := ast.NewIdent(ref.refTable.name)
	refIdent := ast.NewIdent("ref")

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
		Name: joinIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{refIdent},
						Type: &ast.StarExpr{
							X: ref.refTable.structIdent,
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: relationContext(&ast.StarExpr{
							X: jt.structIdent,
						}, &ast.StarExpr{
							X: ref.refTable.structIdent,
						}, &ast.StarExpr{
							X: ref.joinedTable.structIdent,
						}),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: newRelationContext(&ast.StarExpr{
								X: jt.structIdent,
							}, &ast.StarExpr{
								X: ref.refTable.structIdent,
							}, &ast.StarExpr{
								X: ref.joinedTable.structIdent,
							}),
							Args: []ast.Expr{
								jt.recvIdent,
								refIdent,
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
