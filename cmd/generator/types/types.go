package types

import (
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
	Table       *Table
	JoinedTable *JoinedTable
}

type RefJoinedTable struct {
	Table       *JoinedTable
	JoinedTable *JoinedTable
}

type Column struct {
	Name      string
	FieldName string
	Type      ast.Expr
}
