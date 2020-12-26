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

// Reference 外部キーの構造体(yaml用)
type Reference struct {
	Table    string         `yaml:"table"`
	Columns  []*IndexColumn `yaml:"columns"`
	Match    Match          `yaml:"match"`
	OnDelete string         `yaml:"on_delete"`
	OnUpdate string         `yaml:"on_update"`
}
