package codegen

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mazrean/genorm/cmd/generator/types"
)

const (
	genormImport = `"github.com/mazrean/genorm"`
	genormRelationImport = `"github.com/mazrean/genorm/relation"`
	fmtImport    = `"fmt"`
)

var (
	genormIdent = ast.NewIdent("genorm")
	genormRelationIdent = ast.NewIdent("relation")
	fmtIdent    = ast.NewIdent("fmt")
)

func Codegen(packageName string, modulePath string, destinationDir string, baseAst *ast.File, tables []*types.Table, joinedTables []*types.JoinedTable) error {
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
	haveGenormRelation := false
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
					haveGenorm = true
				case genormRelationImport:
					genormRelationIdent = importSpec.Name
					haveGenormRelation = true
				case fmtImport:
					fmtIdent = importSpec.Name
					haveFmt = true
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

		if !haveGenormRelation {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Name: genormRelationIdent,
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: genormRelationImport,
				},
			})
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

func codegenMain(dir *directory, importDecls []ast.Decl, tables []*types.Table, joinedTables []*types.JoinedTable) error {
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

	joinedTableDecls := codegenMainJoinedTableDecls(joinedTables)
	if err != nil {
		return fmt.Errorf("failed to codegen joined tables: %w", err)
	}
	astFile.Decls = append(astFile.Decls, joinedTableDecls...)

	return nil
}

func codegenMainTableDecls(tables []*types.Table) ([]ast.Decl, error) {
	tableDecls := []ast.Decl{}
	for _, table := range tables {
		tableDecls = append(tableDecls, tableDecl(table)...)
	}

	return tableDecls, nil
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

func codegenMainJoinedTableDecls(joinedTables []*types.JoinedTable) []ast.Decl {
	return nil
}
