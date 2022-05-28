package grpcio

import (
	"io"
	"sync"

	"github.com/frantjc/sequence/api/types"
)

type LogStreamServer interface {
	Send(*types.Log) error
}

func NewLogOutStreamWriter(s LogStreamServer) io.Writer {
	return &logStreamWriter{sync.Mutex{}, s, 0}
}

func NewLogErrStreamWriter(s LogStreamServer) io.Writer {
	return &logStreamWriter{sync.Mutex{}, s, 1}
}

func NewLogStreamMultiplexWriter(s LogStreamServer) (io.Writer, io.Writer) {
	return NewLogOutStreamWriter(s), NewLogErrStreamWriter(s)
}

type logStreamWriter struct {
	sync.Mutex
	s      LogStreamServer
	stream int32
}

func (w *logStreamWriter) Write(p []byte) (int, error) {
	if p == nil {
		return len(p), nil
	}

	w.Lock()
	defer w.Unlock()
	err := w.s.Send(&types.Log{
		Data:   p,
		Stream: w.stream,
	})
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

var _ io.Writer = &logStreamWriter{}
