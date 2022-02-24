package genorm

type Table interface {
	Expr
	SQLTableName() string
	Columns() []Column
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
