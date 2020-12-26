package models

// IndexType インデックスの種類
type IndexType string

const (
	// Btree B-treeインデックス
	Btree IndexType = "btree"
	// Hash Hashインデックス
	Hash IndexType = "hash"
)

// IndexDirection インデックスの方向
type IndexDirection string

const (
	// Asc 昇順
	Asc IndexDirection = "asc"
	// Desc 降順
	Desc IndexDirection = "desc"
)

// IndexOption インデックスのオプション
type IndexOption struct {
	KeyBlockSize int    `yaml:"key_block_size"`
	Type         Type   `yaml:"type"`
	Parser       string `yaml:"parser"`
}

// IndexColumn インデックスを貼るカラム
type IndexColumn struct {
	Name      string `yaml:"name"`
	Length    int    `yaml:"length"`
	Direction string `yaml:"direction"`
}
