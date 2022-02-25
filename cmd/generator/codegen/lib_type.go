package codegen

import (
	"go/ast"
)

var (
	columnInterfaceTypeExpr = &ast.SelectorExpr{
		X:   genormIdent,
		Sel: ast.NewIdent("Column"),
	}
	exprTypeInterfaceTypeExpr = &ast.SelectorExpr{
		X:   genormIdent,
		Sel: ast.NewIdent("ExprType"),
	}

	exprExprIdent           = ast.NewIdent("Expr")
	tableExprTableExprIdent = ast.NewIdent("TableExpr")
	typedExprTypedExprIdent = ast.NewIdent("TypedExpr")

	tableColumnsIdent        = ast.NewIdent("Columns")
	tableGetErrors           = ast.NewIdent("GetErrors")
	tableColumnMapIdent      = ast.NewIdent("ColumnMap")
	basicTableTableNameIdent = ast.NewIdent("TableName")

	columnSQLColumnsIdent = ast.NewIdent("SQLColumnName")
	columnTableNameIdent  = ast.NewIdent("TableName")
	columnColumnNameIdent = ast.NewIdent("ColumnName")
)

func wrappedPrimitive(primitive ast.Expr) ast.Expr {
	return &ast.IndexExpr{
		X: &ast.SelectorExpr{
			X:   genormIdent,
			Sel: ast.NewIdent("WrappedPrimitive"),
		},
		Index: primitive,
	}
}
