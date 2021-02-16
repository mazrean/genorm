package models

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Reference 外部キーの構造体(yaml用)
type Reference struct {
	*domain.Reference `yaml:",inline"`
	Columns           []*domain.IndexColumn `yaml:"columns"`
}
