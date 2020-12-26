package models

// Match マッチの種類
type Match string

const (
	// Full 非NULL
	Full Match = "full"
	// Partial 未実装
	Partial Match = "partial"
	// Simple NULL許容
	Simple Match = "simple"
)

// ReferenceOption 外部キーの削除時の処理
type ReferenceOption string

const (
	// Restrict エラーになる
	Restrict ReferenceOption = "restrict"
	// Cascade 参照先の変更に追従する
	Cascade ReferenceOption = "cascade"
	// SetNull NULLに置き換わる
	SetNull ReferenceOption = "set_null"
	// NoAction エラーになる
	NoAction ReferenceOption = "no_action"
)

// ColumnFormat データの保存形式
type ColumnFormat string

const (
	// Fixed 固定長
	Fixed ColumnFormat = "fixed"
	// Dynamic 可変長
	Dynamic ColumnFormat = "dynamic"
)

// Storage カラムの保存先
type Storage string

const (
	// Disk ディスクに格納
	Disk Storage = "disc"
	// Memory メモリに格納
	Memory Storage = "memory"
)

// Column カラムの構造体(yaml用)
type Column struct {
	ID          string       `yaml:"id"`
	Description string       `yaml:"description"`
	Name        string       `yaml:"name"`
	Type        string       `yaml:"type"`
	Null        bool         `yaml:"null"`
	Default     string       `yaml:"dafault"`
	Reference   []*Reference `yaml:"reference"`
	Extra       *Extra       `yaml:"extra"`
}

// Reference 外部キーの構造体(yaml用)
type Reference struct {
	Table    string         `yaml:"table"`
	Columns  []*IndexColumn `yaml:"columns"`
	Match    Match          `yaml:"match"`
	OnDelete string         `yaml:"on_delete"`
	OnUpdate string         `yaml:"on_update"`
}

// Extra 特殊な設定
type Extra struct {
	AutoIncrement bool   `yaml:"auto_increment"`
	Unique        bool   `yaml:"unique"`
	Format        string `yaml:"format"`
	Storage       string `yaml:"strage"`
}
