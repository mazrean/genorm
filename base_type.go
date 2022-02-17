package genorm

import "time"

// Table types

type TableBase interface {
	/*
		TableID
		BasicTable: `table_name`
		JoinTable: ex) `table_name1` JOIN `table_name2` ON `table_name1`.`id` = `table_name2`.`id`
	*/
	TableID() string
	Columns() []Column
	// ColumnMap key: `table_name`.`column_name`
	ColumnMap() map[string]ColumnField
}

type BasicTable interface {
	TableBase
	// TableName `table_name`
	TableName() string
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
	Errors() []error
}

type NullableExpr interface {
	Expr
	IsNull() bool
}

type TableExpr[T TableBase] interface {
	Expr
	TableExpr(T) (string, []any)
}

type TypedExpr[T ExprType] interface {
	Expr
	TypedExpr(T) (string, []any)
}

type TypedTableExpr[T TableBase, S ExprType] interface {
	Expr
	TableExpr[T]
	TypedExpr[S]
}

type NullableTypedTableExpr[T TableBase, S ExprType] interface {
	NullableExpr
	TableExpr[T]
	TypedExpr[S]
}

type ExprType interface {
	bool |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string | time.Time
}

// Column types

type Column interface {
	Expr
	// ColumnID `table_name`.`column_name`
	ColumnID() string
	// TableName `table_name`
	TableName() string
	// ColumnName `column_name`
	ColumnName() string
}

type TableColumns[T TableBase] interface {
	Column
	TableExpr[T]
}

type TypedColumns[S ExprType] interface {
	Column
	TypedExpr[S]
}

type TypedTableColumns[T TableBase, S ExprType] interface {
	TableColumns[T]
	TypedColumns[S]
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
