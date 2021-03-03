package config

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Table 設定ファイルのtableのinterface
type Table interface {
	GetAll() []*domain.Table
	GetPrimaryKeyNames(tableID string) ([]string, error)
	GetColumns(tableID string) ([]*domain.Column, error)
	GetReference(tableID string) ([]*TableReference, error)
}

// TableReference テーブルの外部キー制約
type TableReference struct {
	Column          *domain.Column
	ReferenceTable  *domain.Table
	ReferenceColumn *domain.Column
}
