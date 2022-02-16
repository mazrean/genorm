package genorm

import "time"

// Table types

type TableBase interface {
	TableName() string
}

type BasicTable interface {
	TableBase
	ColumnNames() []string
}

type JoinedTable interface {
	TableBase
	BaseTables() []BasicTable
	Relations() []Relation
	SetRelation(base Column, ref Column, relationType RelationType)
}

// Expr types

type Expr interface {
	Expr() (string, []any)
}

type TableExpr[T TableBase] interface {
	Expr
	TableExpr(T) (string, []any)
}

type TypedExpr[T ColumnType] interface {
	Expr
	TypedExpr(T) (string, []any)
}

type TypedTableExpr[T TableBase, S ColumnType] interface {
	Expr
	TableExpr[T]
	TypedExpr[S]
}

// Column types

type Column interface {
	Expr
	TableName() string
	ColumnName() string
}

type TableColumns[T TableBase] interface {
	Column
	TableExpr[T]
}

type TypedColumns[S ColumnType] interface {
	Column
	TypedExpr[S]
}

type TypedTableColumns[T TableBase, S ColumnType] interface {
	TableColumns[T]
	TypedColumns[S]
}

type ColumnType interface {
	bool |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string | time.Time | []byte
}

// Relation types

type RelationType interface {
	RelationTypeName() string
}

type Join struct{}

func (Join) RelationTypeName() string {
	return "JOIN"
}

type LeftJoin struct{}

func (LeftJoin) RelationTypeName() string {
	return "LEFT JOIN"
}

type RightJoin struct{}

func (RightJoin) RelationTypeName() string {
	return "RIGHT JOIN"
}

type Relation struct {
	base Column
	ref  Column
	Type RelationType
}
