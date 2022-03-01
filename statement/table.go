package statement

import (
	"github.com/mazrean/genorm"
)

type Table interface {
	genorm.Table
	// ColumnMap key: `table_name`.`column_name`
	ColumnMap() map[string]genorm.ColumnFieldExprType
}

type BasicTable interface {
	Table
	genorm.BasicTable
}

type JoinedTable interface {
	Table
	genorm.JoinedTable
}
