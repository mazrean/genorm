package genorm

type Table interface {
	TableName() string
}

type JoinedTable interface {
	Table
	BaseTables() []Table
	Relations() []Relation
	SetRelation(base Column, ref Column, relationType RelationType)
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
