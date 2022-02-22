package generator

import (
	"fmt"
	"go/ast"
	"reflect"
)

type ParserTable struct {
	StructName string
	Columns    []*ParserColumn
	Methods    []*ParserMethod
	RefTables  []*ParserRefTable
}

type ParserMethod struct {
	StructName string
	Type       methodType
	Decl       *ast.FuncDecl
}

type ParserRefTable struct {
	FieldName  string
	StructName string
}

type ParserColumn struct {
	Name      string
	FieldName string
	Type      ast.Expr
}

func parse(f *ast.File) ([]*ParserTable, error) {
	tables := []*ParserTable{}
	methodMap := make(map[string][]*ParserMethod)

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			method, isMethod, err := parseFuncDecl(funcDecl)
			if err != nil {
				return nil, fmt.Errorf("parse func: %w", err)
			}

			if isMethod {
				methodMap[method.StructName] = append(methodMap[method.StructName], method)
			}
		}

		newTables, err := parseGenDecl(genDecl)
		if err != nil {
			return nil, fmt.Errorf("parse gen: %w", err)
		}

		tables = append(tables, newTables...)
	}

	for _, table := range tables {
		table.Methods = methodMap[table.StructName]
	}

	return tables, nil
}

func parseFuncDecl(f *ast.FuncDecl) (*ParserMethod, bool, error) {
	recv := f.Recv
	if recv == nil {
		return nil, false, nil
	}

	if len(recv.List) == 0 {
		return nil, false, nil
	}

	recvType := recv.List[0].Type
	identType, ok := recvType.(*ast.Ident)
	if !ok {
		starType, ok := recvType.(*ast.StarExpr)
		if !ok {
			return nil, false, nil
		}

		identType, ok = starType.X.(*ast.Ident)
		if !ok {
			return nil, false, nil
		}

		return &ParserMethod{
			StructName: identType.Name,
			Type:       methodTypeStar,
			Decl:       f,
		}, true, nil
	}

	return &ParserMethod{
		StructName: identType.Name,
		Type:       methodTypeIdentifier,
		Decl:       f,
	}, true, nil
}

func parseGenDecl(g *ast.GenDecl) ([]*ParserTable, error) {
	tables := []*ParserTable{}

	for _, spec := range g.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok || typeSpec == nil {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok || structType == nil {
			continue
		}

		table, err := parseStructType(typeSpec.Name.Name, structType)
		if err != nil {
			return nil, fmt.Errorf("parse struct: %w", err)
		}

		if table != nil {
			tables = append(tables, table)
		}
	}

	return tables, nil
}

func parseStructType(name string, s *ast.StructType) (*ParserTable, error) {
	fieldList := s.Fields
	if fieldList == nil {
		return nil, nil
	}

	fields := fieldList.List
	if len(fields) == 0 {
		return nil, nil
	}

	columns := []*ParserColumn{}
	refTables := []*ParserRefTable{}
	for _, field := range fields {
		if tableName, isRef := checkRefType(field.Type); isRef {
			for _, name := range field.Names {
				if name == nil {
					continue
				}

				refTables = append(refTables, &ParserRefTable{
					StructName: tableName,
					FieldName:  name.Name,
				})
			}

			continue
		}

		tagLit := field.Tag

		var tag string
		if tagLit != nil {
			structTag := reflect.StructTag(tagLit.Value)
			tag = structTag.Get("genorm")
		}

		for _, name := range field.Names {
			var columnName string
			if len(tag) != 0 {
				columnName = tag
			} else {
				columnName = name.Name
			}

			columns = append(columns, &ParserColumn{
				Name:      columnName,
				FieldName: name.Name,
				Type:      field.Type,
			})
		}
	}

	return &ParserTable{
		StructName: name,
		Columns:    columns,
		RefTables:  refTables,
	}, nil
}

func checkRefType(t ast.Expr) (string, bool) {
	indexExpr, ok := t.(*ast.IndexExpr)
	if !ok || indexExpr == nil {
		return "", false
	}

	selectorExpr, ok := indexExpr.X.(*ast.SelectorExpr)
	if !ok || selectorExpr == nil {
		return "", false
	}

	if selectorExpr.X == nil || selectorExpr.Sel == nil {
		return "", false
	}

	selectorIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return "", false
	}

	if selectorIdent.Name != "genorm" {
		return "", false
	}

	if selectorExpr.Sel.Name != "Ref" {
		return "", false
	}

	ident, ok := indexExpr.Index.(*ast.Ident)
	if !ok || ident == nil {
		return "", false
	}

	return ident.Name, true
}
