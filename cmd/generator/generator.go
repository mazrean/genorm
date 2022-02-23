package generator

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
)

type Config struct {
	JoinNum int
}

func Generate(packageName string, moduleName string, destinationDir string, src io.Reader, config Config) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.Mode(0))
	if err != nil {
		return fmt.Errorf("parse source: %w", err)
	}

	parserTables, err := parse(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	_, _, err = convert(parserTables, config.JoinNum)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	return nil
}
