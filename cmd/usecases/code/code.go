package code

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

import "io"

// GenCode コード生成のinterface
type GenCode interface {
	Write(io.Writer) error
}
