package generator

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"
	"path/filepath"
)

const (
	genormImport = `"github.com/mazrean/genorm"`
	fmtImport    = `"fmt"`
)

var (
	genormIdent = ast.NewIdent("genorm")
	fmtIdent = ast.NewIdent("fmt")
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
