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
		Sel: ast.NewIdent("ColumnFieldExprType"),
	}
	tableTypeExpr = &ast.SelectorExpr{
		X:   genormIdent,
		Sel: ast.NewIdent("Table"),
	}
	basicTableTypeExpr = &ast.SelectorExpr{
		X:   genormIdent,
		Sel: ast.NewIdent("BasicTable"),
	}
	relationTypeExpr = &ast.SelectorExpr{
		X:   genormRelationIdent,
		Sel: ast.NewIdent("Relation"),
	}

	exprExprIdent           = ast.NewIdent("Expr")
	tableExprTableExprIdent = ast.NewIdent("TableExpr")
	typedExprTypedExprIdent = ast.NewIdent("TypedExpr")

	tableColumnsIdent           = ast.NewIdent("Columns")
	tableGetErrorsIdent         = ast.NewIdent("GetErrors")
	tableAddErrorIdent          = ast.NewIdent("AddError")
	tableColumnMapIdent         = ast.NewIdent("ColumnMap")
	basicTableTableNameIdent    = ast.NewIdent("TableName")
	joinedTableBaseTablesIdent  = ast.NewIdent("BaseTables")
	joinedTableSetRelationIdent = ast.NewIdent("SetRelation")

	columnSQLColumnsIdent = ast.NewIdent("SQLColumnName")
	columnTableNameIdent  = ast.NewIdent("TableName")
	columnColumnNameIdent = ast.NewIdent("ColumnName")

	relationJoinedTableNameIdent = ast.NewIdent("JoinedTableName")
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

func typedTableExpr(tableType ast.Expr, exprType ast.Expr) ast.Expr {
	return &ast.IndexListExpr{
		X: &ast.SelectorExpr{
			X:   genormIdent,
			Sel: ast.NewIdent("TypedTableExpr"),
		},
		Indices: []ast.Expr{
			tableType,
			exprType,
		},
	}
}

func typedTableColumn(tableType ast.Expr, exprType ast.Expr) ast.Expr {
	return &ast.IndexListExpr{
		X: &ast.SelectorExpr{
			X:   genormIdent,
			Sel: ast.NewIdent("TypedTableColumns"),
		},
		Indices: []ast.Expr{
			tableType,
			exprType,
		},
	}
}

func relationContext(baseTable ast.Expr, refTable ast.Expr, joinedTable ast.Expr) ast.Expr {
	return &ast.StarExpr{
		X: &ast.IndexListExpr{
			X: &ast.SelectorExpr{
				X:   genormRelationIdent,
				Sel: ast.NewIdent("RelationContext"),
			},
			Indices: []ast.Expr{
				baseTable,
				refTable,
				&ast.StarExpr{
					X: joinedTable,
				},
				joinedTable,
			},
		},
	}
}

func newRelationContext(baseTable ast.Expr, refTable ast.Expr, joinedTable ast.Expr) ast.Expr {
	return &ast.IndexListExpr{
		X: &ast.SelectorExpr{
			X:   genormRelationIdent,
			Sel: ast.NewIdent("NewRelationContext"),
		},
		Indices: []ast.Expr{
			baseTable,
			refTable,
			&ast.StarExpr{
				X: joinedTable,
			},
			joinedTable,
		},
	}
}
