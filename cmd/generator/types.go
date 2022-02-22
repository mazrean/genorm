package generator

import "go/ast"

type Table struct {
	StructName      string
	Columns         []*Column
	Methods         []*Method
	RefTables       []*RefTable
	RefJoinedTables []*RefJoinedTable
}

type Method struct {
	Type methodType
	Decl *ast.FuncDecl
}

type JoinedTable struct {
	Tables          []*Table
	RefTables       []*RefTable
	RefJoinedTables []*RefJoinedTable
}

type methodType int8

const (
	methodTypeIdentifier methodType = iota + 1
	methodTypeStar
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
