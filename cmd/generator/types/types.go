package types

import (
	"errors"
	"go/ast"
)

type Table struct {
	StructName      string
	Columns         []*Column
	Methods         []*Method
	RefTables       []*RefTable
	RefJoinedTables []*RefJoinedTable
}

type Method struct {
	Type MethodType
	Decl *ast.FuncDecl
}

func (m *Method) SetStructName(structName string) error {
	if m.Decl == nil ||
		m.Decl.Recv == nil ||
		len(m.Decl.Recv.List) == 0 ||
		m.Decl.Recv.List[0] == nil ||
		m.Decl.Recv.List[0].Type == nil {
		return errors.New("invalid method")
	}
	switch m.Type {
	case MethodTypeIdentifier:
		ident, ok := m.Decl.Recv.List[0].Type.(*ast.Ident)
		if !ok || ident == nil {
			return errors.New("invalid method")
		}

		ident.Name = structName
	case MethodTypeStar:
		star, ok := m.Decl.Recv.List[0].Type.(*ast.StarExpr)
		if !ok || star == nil {
			return errors.New("invalid method")
		}

		ident, ok := star.X.(*ast.Ident)
		if !ok || ident == nil {
			return errors.New("invalid method")
		}

		ident.Name = structName
	default:
		return errors.New("unknown method type")
	}

	return nil
}

type JoinedTable struct {
	Tables          []*Table
	RefTables       []*RefTable
	RefJoinedTables []*RefJoinedTable
}

type MethodType int8

const (
	MethodTypeIdentifier MethodType = iota + 1
	MethodTypeStar
)

type RefTable struct {
	Table *Table
}

type RefJoinedTable struct {
	Table *JoinedTable
}

type Column struct {
	Name      string
	FieldName string
	Type      ast.Expr
}
