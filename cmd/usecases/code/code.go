package code

import "io"

// GenCode コード生成のinterface
type GenCode interface {
	Write(io.Writer) error
}
