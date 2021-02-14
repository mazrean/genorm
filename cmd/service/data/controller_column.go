package data

// Storage カラムの保存先
type Storage int

const (
	// Disk ディスクに格納
	Disk Storage = iota
	// Memory メモリに格納
	Memory
)

// ColumnFormat データの保存形式
type ColumnFormat int

const (
	// Fixed 固定長
	Fixed ColumnFormat = iota
	// Dynamic 可変長
	Dynamic
)

// Column カラムの構造体(yaml用)
type Column struct {
	ID          string
	Description string
	Name        string
	Type        string
	Null        bool
	Default     string
	Reference   []*Reference
	Extra       *Extra
}

// Reference 外部キーの構造体(yaml用)
type Reference struct {
	Table    string
	Columns  []*IndexColumn
	Match    Match
	OnDelete string
	OnUpdate string
}

// Extra 特殊な設定
type Extra struct {
	AutoIncrement bool
	Unique        bool
	Format        string
	Storage       string
}
