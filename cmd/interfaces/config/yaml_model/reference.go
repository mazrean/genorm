package models

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Reference 外部キーの構造体(yaml用)
type Reference struct {
	Table    string                `yaml:"table"`
	Columns  []*domain.IndexColumn `yaml:"columns"`
	Match    domain.Match          `yaml:"match"`
	OnDelete string                `yaml:"on_delete"`
	OnUpdate string                `yaml:"on_update"`
}
