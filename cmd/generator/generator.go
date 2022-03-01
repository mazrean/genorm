package generator

import (
	"fmt"
	p "go/parser"
	"go/token"
	"io"

	"github.com/mazrean/genorm/cmd/generator/codegen"
	"github.com/mazrean/genorm/cmd/generator/convert"
	"github.com/mazrean/genorm/cmd/generator/parser"
)

type Config struct {
	JoinNum int
}

func Generate(packageName string, moduleName string, destinationDir string, src io.Reader, config Config) error {
	fset := token.NewFileSet()
	f, err := p.ParseFile(fset, "", src, p.Mode(0))
	if err != nil {
		return fmt.Errorf("parse source: %w", err)
	}

	parserTables, err := parser.Parse(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	tables, joinedTables, err := convert.Convert(parserTables, config.JoinNum)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	err = codegen.Codegen(packageName, moduleName, destinationDir, f, tables, joinedTables)
	if err != nil {
		return fmt.Errorf("codegen: %w", err)
	}

	return nil
}
