package domain

// ColumnFormat データの保存形式
type ColumnFormat string

const (
	// Fixed 固定長
	Fixed ColumnFormat = "fixed"
	// Dynamic 可変長
	Dynamic ColumnFormat = "dynamic"
)

// Storage カラムの保存先
type Storage string

const (
	// Disk ディスクに格納
	Disk Storage = "disc"
	// Memory メモリに格納
	Memory Storage = "memory"
)

// Column カラムの構造体(yaml用)
type Column struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Null        bool   `yaml:"null"`
	Default     string `yaml:"dafault"`
}

// Extra 特殊な設定
type Extra struct {
	AutoIncrement bool   `yaml:"auto_increment"`
	Unique        bool   `yaml:"unique"`
	Format        string `yaml:"format"`
	Storage       string `yaml:"strage"`
}
