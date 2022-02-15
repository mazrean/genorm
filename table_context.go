package genorm

type TableContext[T TableBase] struct {
	*Context[T]
}

func Table[T TableBase](table T) *TableContext[T] {
	return &TableContext[T]{
		Context: newContext(table),
	}
}

func (tc *TableContext[TableBase]) Create(tableBases ...TableBase) *CreateContext[TableBase] {
	return newCreateContext(tc.Context, tableBases...)
}
