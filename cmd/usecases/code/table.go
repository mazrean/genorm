package code

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

import (
	"context"
	"io"

	"github.com/mazrean/gopendb-generator/cmd/domain"
)

// Table コード生成のtableのinterface
type Table interface {
	Generate(ctx context.Context, tableDetail *TableDetail, fileWriter func(path string) (io.WriteCloser, error)) error
}

// TableDetail テーブルの詳細
type TableDetail struct {
	*domain.Table
	PrimaryKeyColumnIDs []string
	Columns             []*domain.Column
	References          []*TableReference
}

// TableReference テーブルの外部キー制約
type TableReference struct {
	Column          *domain.Column
	ReferenceTable  *domain.Table
	ReferenceColumn *domain.Column
}
