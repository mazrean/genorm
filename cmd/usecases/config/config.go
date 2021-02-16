package config

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

import (
	"github.com/mazrean/gopendb-generator/cmd/domain"
)

// Config configの取得
type Config interface {
	Get() (*domain.Config, error)
}
