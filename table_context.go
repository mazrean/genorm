package genorm

type TableContext[T TableBase] struct {
	table T
}

func Table[T TableBase](table T) *TableContext[T] {
	return &TableContext[T]{
		table: table,
	}
}

func (tc *TableContext[TableBase]) Create(tableBases ...TableBase) *CreateContext[TableBase] {
	return newCreateContext(tc.table, tableBases...)
}
