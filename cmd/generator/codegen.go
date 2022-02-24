package generator

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	genormImport = `"github.com/mazrean/genorm"`
	fmtImport    = `"fmt"`
)

var (
	genormIdent = ast.NewIdent("genorm")
	fmtIdent    = ast.NewIdent("fmt")
)

func codegen(packageName string, modulePath string, destinationDir string, baseAst *ast.File, tables []*Table, joinedTables []*JoinedTable) error {
	dir, err := newDirectory(destinationDir, packageName, modulePath)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	importDecls := codegenImportDecls(baseAst)

	return nil
}

func codegenImportDecls(baseAst *ast.File) []ast.Decl {
	importDecls := []ast.Decl{}
	haveGenorm := false
	haveFmt := false
	for _, decl := range baseAst.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl == nil || genDecl.Tok != token.IMPORT || len(genDecl.Specs) == 0 {
			continue
		}

		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok || importSpec == nil {
				continue
			}

			if importSpec.Name != nil {
				switch importSpec.Path.Value {
				case genormImport:
					genormIdent = importSpec.Name
				case fmtImport:
					fmtIdent = importSpec.Name
				}
			}
		}

		if !haveGenorm {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Name: genormIdent,
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: genormImport,
				},
			})

			haveGenorm = true
		}

		if !haveFmt {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Name: fmtIdent,
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmtImport,
				},
			})

			haveFmt = true
		}

		importDecls = append(importDecls, genDecl)
	}

	return importDecls
}

func codegenMain(dir *directory, importDecls []ast.Decl, tables []*Table, joinedTables []*JoinedTable) error {
	f, err := dir.addFile("genorm.go")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	astFile := f.ast()

	astFile.Decls = append(astFile.Decls, importDecls...)

	tableDecls, err := codegenMainTableDecls(tables)
	if err != nil {
		return fmt.Errorf("failed to codegen tables: %w", err)
	}
	astFile.Decls = append(astFile.Decls, tableDecls...)

	joinedTableDecls, err := codegenMainJoinedTableDecls(joinedTables)
	if err != nil {
		return fmt.Errorf("failed to codegen joined tables: %w", err)
	}
	astFile.Decls = append(astFile.Decls, joinedTableDecls...)

	return nil
}

func codegenMainTableDecls(tables []*Table) ([]ast.Decl, error) {
	tableDecls := []ast.Decl{}
	for _, table := range tables {
		tableDecls = append(tableDecls, codegenMainTableDecl(table)...)
	}

	return tableDecls, nil
}

func codegenMainTableDecl(table *Table) []ast.Decl {
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
		columnExpr, newColumnDecls := codegenMainColumnDecl(table, column)
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
		Name: tableSQLTableNameIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{Type: ast.NewIdent("string")},
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
								Sel: basicTableTableNameIdent,
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
	})

	tableDecls = append(tableDecls, columnDecls...)

	return tableDecls
}

func codegenFieldTypeExpr(columnTypeExpr ast.Expr) ast.Expr {
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

func codegenMainColumnDecl(table *Table, column *Column) (ast.Expr, []ast.Decl) {
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
		})
	}

	return columnVarIdent, columnDecls
}

type directory struct {
	path        string
	packageName string
	modulePath  string
}

func newDirectory(path string, packageName string, modulePath string) (*directory, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	return &directory{
		path:        path,
		packageName: packageName,
		modulePath:  modulePath,
	}, nil
}

func (d *directory) addDirectory(name string, packageName string) (*directory, error) {
	return newDirectory(filepath.Join(d.path, name), packageName, filepath.Join(d.modulePath, name))
}

func (d *directory) addFile(name string) (*file, error) {
	os.Create(filepath.Join(d.path, name))

	return nil, nil
}

type file struct {
	writer io.WriteCloser
	file   *ast.File
}

func newFile(path string, packageName string) (*file, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	astFile := &ast.File{
		Name: ast.NewIdent(packageName),
	}

	return &file{
		writer: f,
		file:   astFile,
	}, nil
}

func (f *file) ast() *ast.File {
	return f.file
}

func (f *file) Close() (err error) {
	defer f.writer.Close()

	err = format.Node(f.writer, token.NewFileSet(), f.file)
	if err != nil {
		return fmt.Errorf("failed to format file: %w", err)
	}

	return nil
}
