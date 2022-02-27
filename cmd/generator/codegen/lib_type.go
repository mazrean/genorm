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
	columnFieldExprTypeExpr = &ast.SelectorExpr{
		X:   genormIdent,
		Sel: ast.NewIdent("ColumnFieldExpr"),
	}
	basicTableTypeExpr = &ast.SelectorExpr{
		X:   genormIdent,
		Sel: ast.NewIdent("BasicTable"),
	}

	exprExprIdent           = ast.NewIdent("Expr")
	tableExprTableExprIdent = ast.NewIdent("TableExpr")
	typedExprTypedExprIdent = ast.NewIdent("TypedExpr")

	tableColumnsIdent          = ast.NewIdent("Columns")
	tableGetErrorsIdent        = ast.NewIdent("GetErrors")
	tableAddErrorIdent         = ast.NewIdent("AddError")
	tableColumnMapIdent        = ast.NewIdent("ColumnMap")
	basicTableTableNameIdent   = ast.NewIdent("TableName")
	joinedTableBaseTablesIdent = ast.NewIdent("BaseTables")

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

func typedTableColumn(tableType ast.Expr, exprType ast.Expr) ast.Expr {
	return &ast.SelectorExpr{
		X:   genormIdent,
		Sel: ast.NewIdent("TypedTableColumn"),
	}
}
