package writer

import "io"

// Writer io.Writerの生成のinterface
type Writer interface {
	FileWriter(file string) (io.Writer, error)
}
