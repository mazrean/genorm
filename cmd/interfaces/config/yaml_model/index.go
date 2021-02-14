package models

import "github.com/mazrean/gopendb-generator/cmd/domain"

// IndexOption インデックスのオプション
type IndexOption struct {
	KeyBlockSize int `yaml:"key_block_size"`
	*domain.Type `yaml:"type"`
	Parser       string `yaml:"parser"`
}
