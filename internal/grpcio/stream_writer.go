package grpcio

import (
	"io"

	"github.com/frantjc/sequence/api/types"
)

func NewLogStreamWriter(s LogStream) io.Writer {
	return &logStreamWriter{s}
}

type LogStream interface {
	Send(*types.Log) error
}

type logStreamWriter struct {
	s LogStream
}

var _ io.Writer = &logStreamWriter{}

func (w *logStreamWriter) Write(p []byte) (int, error) {
	err := w.s.Send(&types.Log{
		Line: string(p),
	})
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
