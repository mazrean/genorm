package genorm

// ExportSubQuery exports SubQuery function for testing in external packages
func ExportSubQuery[TOuter Table, TInner Table, S ExprType](
	pluckCtx *PluckContext[TInner, S],
) TypedTableExpr[TOuter, S] {
	return SubQuery[TOuter](pluckCtx)
}
