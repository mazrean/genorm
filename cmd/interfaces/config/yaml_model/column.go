package models

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Column カラムの構造体(yaml用)
type Column struct {
	*domain.Column `yaml:",inline"`
	Reference      []*Reference  `yaml:"reference"`
	Extra          *domain.Extra `yaml:"extra"`
}
