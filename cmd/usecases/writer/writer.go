package writer

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

import (
	"context"
	"io"
)

// Writer io.Writerの生成のinterface
type Writer interface {
	FileWriterGenerator(ctx context.Context, progress *Progress, directoryPath string) (func(path string) (io.WriteCloser, error), error)
}

// Progress 進捗状況の伝達
type Progress struct {
	Total    chan<- struct{}
	Progress chan<- struct{}
}
