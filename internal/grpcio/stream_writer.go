package grpcio

import (
	"io"
	"sync"

	"github.com/frantjc/sequence/api/types"
)

func NewLogStreamWriter(s LogStreamServer) io.Writer {
	return &logStreamWriter{sync.Mutex{}, s}
}

type LogStreamServer interface {
	Send(*types.Log) error
}

type logStreamWriter struct {
	sync.Mutex
	s LogStreamServer
}

var _ io.Writer = &logStreamWriter{}

func (w *logStreamWriter) Write(p []byte) (int, error) {
	w.Lock()
	defer w.Unlock()
	err := w.s.Send(&types.Log{
		Out: string(p),
	})
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
