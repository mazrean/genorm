package parser

import (
	"fmt"
	"go/ast"
	"reflect"

	"github.com/mazrean/genorm/cmd/generator/types"
)

type parserTable struct {
	StructName string
	Columns    []*parserColumn
	Methods    []*parserMethod
	RefTables  []*parserRefTable
}

type parserMethod struct {
	StructName string
	Type       types.MethodType
	Decl       *ast.FuncDecl
}

type parserRefTable struct {
	FieldName  string
	StructName string
}

type parserColumn struct {
	Name      string
	FieldName string
	Type      ast.Expr
}

func Parse(f *ast.File) ([]*types.Table, error) {
	parserTables := []*parserTable{}
	methodMap := make(map[string][]*parserMethod)

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

			continue
		}

		newTables, err := parseGenDecl(genDecl)
		if err != nil {
			return nil, fmt.Errorf("parse gen: %w", err)
		}

		parserTables = append(parserTables, newTables...)
	}

	for _, parserTable := range parserTables {
		parserTable.Methods = methodMap[parserTable.StructName]
	}

	tables, err := convertTables(parserTables)
	if err != nil {
		return nil, fmt.Errorf("convert tables: %w", err)
	}

	return tables, nil
}

func convertTables(tables []*parserTable) ([]*types.Table, error) {
	type tablePair struct {
		parser    *parserTable
		converted *types.Table
	}

	pairMap := make(map[string]*tablePair, len(tables))
	for _, table := range tables {
		pairMap[table.StructName] = &tablePair{
			parser:    table,
			converted: convertTable(table),
		}
	}

	convertedTables := make([]*types.Table, 0, len(tables))
	for _, pair := range pairMap {
		refTables := make([]*types.RefTable, 0, len(pair.parser.RefTables))
		for _, refParserTable := range pair.parser.RefTables {
			refTable, ok := pairMap[refParserTable.StructName]
			if !ok {
				return nil, fmt.Errorf("ref table not found: %s", refParserTable.StructName)
			}

			refTables = append(refTables, &types.RefTable{
				Table: refTable.converted,
			})
		}

		pair.converted.RefTables = refTables

		convertedTables = append(convertedTables, pair.converted)
	}

	return convertedTables, nil
}

func convertTable(table *parserTable) *types.Table {
	columns := make([]*types.Column, 0, len(table.Columns))
	for _, column := range table.Columns {
		columns = append(columns, &types.Column{
			Name:      column.Name,
			FieldName: column.FieldName,
			Type:      column.Type,
		})
	}

	methods := make([]*types.Method, 0, len(table.Methods))
	for _, method := range table.Methods {
		methods = append(methods, &types.Method{
			Type: method.Type,
			Decl: method.Decl,
		})
	}

	return &types.Table{
		StructName: table.StructName,
		Columns:    columns,
		Methods:    methods,
	}
}

func parseFuncDecl(f *ast.FuncDecl) (*parserMethod, bool, error) {
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

		return &parserMethod{
			StructName: identType.Name,
			Type:       types.MethodTypeStar,
			Decl:       f,
		}, true, nil
	}

	return &parserMethod{
		StructName: identType.Name,
		Type:       types.MethodTypeIdentifier,
		Decl:       f,
	}, true, nil
}

func parseGenDecl(g *ast.GenDecl) ([]*parserTable, error) {
	tables := []*parserTable{}

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

func parseStructType(name string, s *ast.StructType) (*parserTable, error) {
	fieldList := s.Fields
	if fieldList == nil {
		return nil, nil
	}

	fields := fieldList.List
	if len(fields) == 0 {
		return nil, nil
	}

	columns := []*parserColumn{}
	refTables := []*parserRefTable{}
	for _, field := range fields {
		if tableName, isRef := checkRefType(field.Type); isRef {
			for _, name := range field.Names {
				if name == nil {
					continue
				}

				refTables = append(refTables, &parserRefTable{
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

			columns = append(columns, &parserColumn{
				Name:      columnName,
				FieldName: name.Name,
				Type:      field.Type,
			})
		}
	}

	return &parserTable{
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
