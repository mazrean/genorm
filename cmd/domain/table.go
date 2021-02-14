package domain

// Table テーブルの構造体(yaml用)
type Table struct {
	ID           string `yaml:"id"`
	Description  string `yaml:"description"`
	Name         string `yaml:"name"`
	Engin        string `yaml:"engin"`
	CharSet      string `yaml:"char_set"`
	MaxRows      int64  `yaml:"max_rows"`
	MinRows      int64  `yaml:"min_rows"`
	AvgRowLength int64  `yaml:"avg_row_length"`
}
