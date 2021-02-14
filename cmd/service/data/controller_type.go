package data

// DataType DBの型一覧
type DataType int

const (
	// Bit BIT
	Bit DataType = iota
	// TinyInt TINYINT
	TinyInt
	// SmallInt SMALLINT
	SmallInt
	// MediumInt MEDIUMINT
	MediumInt
	// Int INT
	Int
	// Integer INTEGER
	Integer
	// BigInt BIGINT
	BigInt
	// Real REAL
	Real
	// Double DOUBLE
	Double
	// Float FLOAT
	Float
	// Decimal DECIMAL
	Decimal
	// Numeric NUMERIC
	Numeric
	// Date DATE
	Date
	// Time TIME
	Time
	// Timestamp TIMESTAMP
	Timestamp
	// Datetime DATETIME
	Datetime
	// Year YEAR
	Year
	// Char CHAR
	Char
	// Varchar VARCHAR
	Varchar
	// Binary BINARY
	Binary
	// Varbinary VARBINARY
	Varbinary
	// TinyBlob TINYBLOB
	TinyBlob
	// Blob BLOB
	Blob
	// MediumBlob MEDIUMBLOB
	MediumBlob
	// LongBlob LONGBLOB
	LongBlob
	// TinyText TINYTEXT
	TinyText
	// Text TEXT
	Text
	// MediumText MEDIUMTEXT
	MediumText
	// LongText LONGTEXT
	LongText
	// Enum ENUM
	Enum
	// Set SET
	Set
)

// Type DBの型(yaml用)
type Type struct {
	Name   DataType
	Length int
	Values []string
}
