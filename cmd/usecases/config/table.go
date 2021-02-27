package config

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Table 設定ファイルのtableのinterface
type Table interface {
	GetAll() []*TableDetail
	GetColumns(tableID string) []*domain.Column
	GetReference(tableID string) []*TableReference
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
