package models

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Table テーブルの構造体(yaml用)
type Table struct {
	ID           string      `yaml:"id"`
	Description  string      `yaml:"description"`
	Name         string      `yaml:"name"`
	PrimaryKey   *PrimaryKey `yaml:"primary_key"`
	Engin        string      `yaml:"engin"`
	CharSet      string      `yaml:"char_set"`
	MaxRows      int64       `yaml:"max_rows"`
	MinRows      int64       `yaml:"min_rows"`
	AvgRowLength int64       `yaml:"avg_row_length"`
	Columns      []*Column   `yaml:"columns"`
}

// PrimaryKey 主キーの構造体(yaml)
type PrimaryKey struct {
	Columns []*domain.IndexColumn `yaml:"columns"`
	Options *IndexOption          `yaml:"options"`
}
