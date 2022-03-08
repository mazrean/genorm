package genorm

type Table interface {
	Expr
	New() Table
	Columns() []Column
	// ColumnMap key: `table_name`.`column_name`
	ColumnMap() map[string]ColumnFieldExprType
	GetErrors() []error
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
