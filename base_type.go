package genorm

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

type TableColumns[T TableBase] interface {
	TableColumnName(T) string
}

type Column interface {
	TableName() string
	ColumnName() string
}

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
