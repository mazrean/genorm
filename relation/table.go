package relation

import (
	"github.com/mazrean/genorm"
)

type Table interface {
	genorm.Table
}

type BasicTable interface {
	Table
	genorm.BasicTable
}

type JoinedTable interface {
	Table
	genorm.JoinedTable
	SetRelation(*Relation)
}
