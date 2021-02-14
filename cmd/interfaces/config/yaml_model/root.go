package models

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Root yamlのルート
type Root struct {
	Config *domain.Config `yaml:"config"`
	Tables []*Table       `yaml:"tables"`
}
