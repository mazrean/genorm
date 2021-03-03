package code

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

import (
	"context"

	"github.com/mazrean/gopendb-generator/cmd/domain"
)

// Table コード生成のtableのinterface
type Table interface {
	Generate(ctx context.Context, progress *Progress, tableDetail *TableDetail) error
}

// TableDetail テーブルの詳細
type TableDetail struct {
	*domain.Table
	PrimaryKeyColumnNames []string
	Columns               []*domain.Column
	References            []*TableReference
}

// TableReference テーブルの外部キー制約
type TableReference struct {
	Column          *domain.Column
	ReferenceTable  *domain.Table
	ReferenceColumn *domain.Column
}

// Progress 進捗状況の伝達
type Progress struct {
	Total    chan<- struct{}
	Progress chan<- struct{}
}
