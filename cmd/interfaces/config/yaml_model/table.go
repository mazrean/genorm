package models

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Table テーブルの構造体(yaml用)
type Table struct {
	*domain.Table `yaml:",inline"`
	PrimaryKey    *PrimaryKey `yaml:"primary_key"`
	Columns       []*Column   `yaml:"columns"`
}

// PrimaryKey 主キーの構造体(yaml)
type PrimaryKey struct {
	Columns []*domain.IndexColumn `yaml:"columns"`
	Options *IndexOption          `yaml:"options"`
}
