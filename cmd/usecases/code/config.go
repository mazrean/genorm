package code

import (
	"github.com/mazrean/gopendb-generator/cmd/domain"
)

// Config コード生成のコンフィグ設定のinterface
type Config interface {
	Set(*domain.Config)
}
