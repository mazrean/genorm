package domain

// DataType DBの型一覧
type DataType string

const (
	// Bit BIT
	Bit DataType = "bit"
	// TinyInt TINYINT
	TinyInt = "tiny_int"
	// SmallInt SMALLINT
	SmallInt = "small_int"
	// MediumInt MEDIUMINT
	MediumInt = "medium_int"
	// Int INT
	Int = "int"
	// Integer INTEGER
	Integer = "integer"
	// BigInt BIGINT
	BigInt = "big_int"
	// Real REAL
	Real = "real"
	// Double DOUBLE
	Double = "double"
	// Float FLOAT
	Float = "float"
	// Decimal DECIMAL
	Decimal = "decimal"
	// Numeric NUMERIC
	Numeric = "numeric"
	// Date DATE
	Date = "date"
	// Time TIME
	Time = "time"
	// Timestamp TIMESTAMP
	Timestamp = "timestamp"
	// Datetime DATETIME
	Datetime = "datetime"
	// Year YEAR
	Year = "year"
	// Char CHAR
	Char = "char"
	// Varchar VARCHAR
	Varchar = "varchar"
	// Binary BINARY
	Binary = "binary"
	// Varbinary VARBINARY
	Varbinary = "varbinary"
	// TinyBlob TINYBLOB
	TinyBlob = "tiny_blob"
	// Blob BLOB
	Blob = "blob"
	// MediumBlob MEDIUMBLOB
	MediumBlob = "medium_blob"
	// LongBlob LONGBLOB
	LongBlob = "long_blob"
	// TinyText TINYTEXT
	TinyText = "tiny_text"
	// Text TEXT
	Text = "text"
	// MediumText MEDIUMTEXT
	MediumText = "medium_text"
	// LongText LONGTEXT
	LongText = "long_text"
	// Enum ENUM
	Enum = "enum"
	// Set SET
	Set = "set"
)

// Type DBの型(yaml用)
type Type struct {
	Name   DataType `yaml:"name"`
	Length int      `yaml:"length"`
	Values []string `yaml:"values"`
}
