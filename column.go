package genorm

type Column interface {
	Expr
	// SQLColumnName `table_name`.`column_name`
	SQLColumnName() string
	// TableName table_name
	TableName() string
	// ColumnName column_name
	ColumnName() string
}

type TableColumns[T Table] interface {
	Column
	TableExpr[T]
}

type TypedColumns[S ExprType] interface {
	Column
	TypedExpr[S]
}

type TypedTableColumns[T Table, S ExprType] interface {
	TableColumns[T]
	TypedColumns[S]
}
