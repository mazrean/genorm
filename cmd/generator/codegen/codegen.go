package codegen

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mazrean/genorm/cmd/generator/types"
)

const (
	genormImport          = `"github.com/mazrean/genorm"`
	genormRelationImport  = `"github.com/mazrean/genorm/relation"`
	genormStatementImport = `"github.com/mazrean/genorm/statement"`
	fmtImport             = `"fmt"`
)

var (
	genormIdent          = ast.NewIdent("genorm")
	genormRelationIdent  = ast.NewIdent("relation")
	genormStatementIdent = ast.NewIdent("statement")
	fmtIdent             = ast.NewIdent("fmt")

	rootPackageIdent *ast.Ident
)

func Codegen(
	packageName string,
	modulePath string,
	destinationDir string,
	baseAst *ast.File,
	tables []*types.Table,
	joinedTables []*types.JoinedTable,
) error {
	rootPackageIdent = ast.NewIdent(packageName)

	dir, err := newDirectory(destinationDir, packageName, modulePath)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	importDecls := codegenImportDecls(baseAst)

	codegenTables, codegenJoinedTables, err := convert(tables, joinedTables)
	if err != nil {
		return fmt.Errorf("failed to convert tables: %w", err)
	}

	err = codegenMain(dir, importDecls, codegenTables, codegenJoinedTables)
	if err != nil {
		return fmt.Errorf("failed to codegen main: %w", err)
	}

	for _, table := range codegenTables {
		err = codegenTable(dir, importDecls, table)
		if err != nil {
			return fmt.Errorf("failed to codegen table(%s): %w", table.name, err)
		}
	}

	return nil
}

func codegenImportDecls(baseAst *ast.File) []ast.Decl {
	importDecls := []ast.Decl{}

	haveImport := false
	haveGenorm := false
	haveGenormRelation := false
	haveGenormStatement := false
	haveFmt := false
	for _, decl := range baseAst.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl == nil || genDecl.Tok != token.IMPORT || len(genDecl.Specs) == 0 {
			continue
		}

		haveImport = true

		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok || importSpec == nil {
				continue
			}

			switch importSpec.Path.Value {
			case genormImport:
				if importSpec.Name != nil {
					genormIdent = importSpec.Name
				}
				haveGenorm = true
			case genormRelationImport:
				if importSpec.Name != nil {
					genormRelationIdent = importSpec.Name
				}
				haveGenormRelation = true
			case genormStatementImport:
				if importSpec.Name != nil {
					genormStatementIdent = importSpec.Name
				}
				haveGenormStatement = true
			case fmtImport:
				if importSpec.Name != nil {
					fmtIdent = importSpec.Name
				}
				haveFmt = true
			}
		}

		if !haveGenorm {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: genormImport,
				},
			})

			haveGenorm = true
		}

		if !haveGenormRelation {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: genormRelationImport,
				},
			})
		}

		if !haveGenormStatement {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: genormStatementImport,
				},
			})
		}

		if !haveFmt {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmtImport,
				},
			})

			haveFmt = true
		}

		importDecls = append(importDecls, genDecl)
	}

	if !haveImport {
		importDecls = append(importDecls, &ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Name: genormIdent,
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: genormImport,
					},
				},
				&ast.ImportSpec{
					Name: genormRelationIdent,
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: genormRelationImport,
					},
				},
				&ast.ImportSpec{
					Name: fmtIdent,
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmtImport,
					},
				},
			},
		})
	}

	return importDecls
}

type refTable struct {
	refTable    *table
	joinedTable *joinedTable
}

type refJoinedTable struct {
	refTable    *joinedTable
	joinedTable *joinedTable
}

func convert(tables []*types.Table, joinedTables []*types.JoinedTable) ([]*table, []*joinedTable, error) {
	tableMap := make(map[string]*table, len(tables))
	codegenTables := make([]*table, 0, len(tables))
	for _, table := range tables {
		codegenTable, err := newTable(table)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create table(%s): %w", table.StructName, err)
		}

		tableMap[table.StructName] = codegenTable
		codegenTables = append(codegenTables, codegenTable)
	}

	joinedTableMap := make(map[string]*joinedTable, len(joinedTables))
	codegenJoinedTables := make([]*joinedTable, 0, len(joinedTables))
	for _, joinedTable := range joinedTables {
		codegenJoinedTable := newJoinedTable(joinedTable)

		joinedTableMap[joinedTableName(joinedTable)] = codegenJoinedTable
		codegenJoinedTables = append(codegenJoinedTables, codegenJoinedTable)
	}

	for _, typesTable := range tables {
		codegenTable := tableMap[typesTable.StructName]

		refTables := make([]*refTable, 0, len(typesTable.RefTables))
		for _, typeRefTable := range typesTable.RefTables {
			refTables = append(refTables, &refTable{
				refTable:    tableMap[typeRefTable.Table.StructName],
				joinedTable: joinedTableMap[joinedTableName(typeRefTable.JoinedTable)],
			})
		}
		codegenTable.refTables = refTables

		refJoinedTables := make([]*refJoinedTable, 0, len(typesTable.RefJoinedTables))
		for _, typeRefJoinedTable := range typesTable.RefJoinedTables {
			refJoinedTables = append(refJoinedTables, &refJoinedTable{
				refTable:    joinedTableMap[joinedTableName(typeRefJoinedTable.Table)],
				joinedTable: joinedTableMap[joinedTableName(typeRefJoinedTable.JoinedTable)],
			})
		}
		codegenTable.refJoinedTables = refJoinedTables
	}

	for _, typesJoinedTable := range joinedTables {
		codegenJoinedTable := joinedTableMap[joinedTableName(typesJoinedTable)]

		tables := make([]*table, 0, len(typesJoinedTable.Tables))
		for _, typeTable := range typesJoinedTable.Tables {
			tables = append(tables, tableMap[typeTable.StructName])
		}
		codegenJoinedTable.tables = tables

		refTables := make([]*refTable, 0, len(typesJoinedTable.RefTables))
		for _, typeRefTable := range typesJoinedTable.RefTables {
			refTables = append(refTables, &refTable{
				refTable:    tableMap[typeRefTable.Table.StructName],
				joinedTable: joinedTableMap[joinedTableName(typeRefTable.JoinedTable)],
			})
		}
		codegenJoinedTable.refTables = refTables

		refJoinedTables := make([]*refJoinedTable, 0, len(typesJoinedTable.RefJoinedTables))
		for _, typeRefJoinedTable := range typesJoinedTable.RefJoinedTables {
			refJoinedTables = append(refJoinedTables, &refJoinedTable{
				refTable:    joinedTableMap[joinedTableName(typeRefJoinedTable.Table)],
				joinedTable: joinedTableMap[joinedTableName(typeRefJoinedTable.JoinedTable)],
			})
		}
		codegenJoinedTable.refJoinedTables = refJoinedTables
	}

	return codegenTables, codegenJoinedTables, nil
}

func codegenMain(dir *directory, importDecls []ast.Decl, tables []*table, joinedTables []*joinedTable) error {
	f, err := dir.addFile("genorm.go")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	astFile := f.ast()

	astFile.Decls = append(astFile.Decls, importDecls...)

	for _, table := range tables {
		astFile.Decls = append(astFile.Decls, table.decl()...)
	}

	for _, joinedTable := range joinedTables {
		astFile.Decls = append(astFile.Decls, joinedTable.decl()...)
	}

	return nil
}

func codegenTable(dir *directory, importDecls []ast.Decl, table *table) error {
	rootModulePath := dir.modulePath

	dir, err := dir.addDirectory(table.snakeName(), table.lowerName())
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := dir.addFile(fmt.Sprintf("%s.go", table.snakeName()))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	astFile := f.ast()

	astFile.Decls = append(astFile.Decls, importDecls...)

	astFile.Decls = append(astFile.Decls, &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, rootModulePath),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: genormImport,
				},
			},
		},
	})

	astFile.Decls = append(astFile.Decls, table.tablePackageDecls()...)

	return nil
}
