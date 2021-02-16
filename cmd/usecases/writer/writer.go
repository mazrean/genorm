package writer

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

import "io"

// Writer io.Writerの生成のinterface
type Writer interface {
	FileWriter(file string) (io.Writer, error)
}
