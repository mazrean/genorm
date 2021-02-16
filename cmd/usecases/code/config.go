package code

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

import (
	"github.com/mazrean/gopendb-generator/cmd/domain"
)

// Config コード生成のコンフィグ設定のinterface
type Config interface {
	Set(*domain.Config)
}
