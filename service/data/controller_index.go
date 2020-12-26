package data

// IndexOption インデックスのオプション
type IndexOption struct {
	KeyBlockSize int
	Type         Type
	Parser       string
}

// IndexColumn インデックスを貼るカラム
type IndexColumn struct {
	Name      string
	Length    int
	Direction string
}
