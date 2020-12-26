package data

// Table テーブルの構造体
type Table struct {
	ID           string
	Description  string
	Name         string
	PrimaryKey   *PrimaryKey
	Engin        string
	CharSet      string
	MaxRows      int64
	MinRows      int64
	AvgRowLength int64
	Columns      []*Column
}

// PrimaryKey 主キーの構造体(yaml)
type PrimaryKey struct {
	Columns []*IndexColumn
	Options *IndexOption
}
