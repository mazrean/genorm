package data

// Match マッチの種類
type Match int

const (
	// Full 非NULL
	Full Match = iota
	// Partial 未実装
	Partial
	// Simple NULL許容
	Simple
)

// ReferenceOption 外部キーの削除時の処理
type ReferenceOption int

const (
	// Restrict エラーになる
	Restrict ReferenceOption = iota
	// Cascade 参照先の変更に追従する
	Cascade
	// SetNull NULLに置き換わる
	SetNull
	// NoAction エラーになる
	NoAction
)
