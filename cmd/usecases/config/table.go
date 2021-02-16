package config

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Table 設定ファイルのtableのinterface
type Table interface {
	GetAll() ([]*TableDetail, error)
	GetColumns(tableID string) ([]*domain.Column, error)
	GetReference(tableID string) ([]*TableReference, error)
}

// TableDetail テーブルの詳細
type TableDetail struct {
	*domain.Table
	PrimaryKeyColumnIDs []string
}

// TableReference テーブルの外部キー制約
type TableReference struct {
	Column          *domain.Column
	ReferenceTable  *domain.Table
	ReferenceColumn *domain.Column
}
