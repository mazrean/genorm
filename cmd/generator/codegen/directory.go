package codegen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"
	"path/filepath"
)

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
