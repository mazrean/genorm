package statement

import (
	"database/sql"
	"errors"

	"github.com/mazrean/genorm"
)

type Table interface {
	genorm.Table
	// ColumnMap key: `table_name`.`column_name`
	ColumnMap() map[string]ColumnField
}

type BasicTable interface {
	Table
	genorm.BasicTable
}

type JoinedTable interface {
	Table
	genorm.JoinedTable
}

type ColumnField interface {
	sql.Scanner
	iValue() (val any, err error)
}

var (
	ErrNullValue   = errors.New("null value")
	ErrEmptyColumn = errors.New("empty column")
)
