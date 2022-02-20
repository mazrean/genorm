package generator

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
)

func Generate(packageName string, moduleName string, destinationDir string, src io.Reader) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.Mode(0))
	if err != nil {
		return fmt.Errorf("parse source: %w", err)
	}

	_, err = parse(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	return nil
}
