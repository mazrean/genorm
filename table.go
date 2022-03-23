package genorm

type Table interface {
	Expr
	Columns() []Column
	// ColumnMap key: table_name.column_name
	ColumnMap() map[string]ColumnFieldExprType
	GetErrors() []error
}

type TablePointer[T any] interface {
	Table
	*T
}

type BasicTable interface {
	Table
	// TableName table_name
	TableName() string
}

type JoinedTable interface {
	Table
	BaseTables() []BasicTable
	AddError(error)
}
