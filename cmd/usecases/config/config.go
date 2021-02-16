package config

import (
	"github.com/mazrean/gopendb-generator/cmd/domain"
)

// Config configの取得
type Config interface {
	Get() (*domain.Config, error)
}
