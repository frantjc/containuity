package grpcio

import (
	"io"
	"sync"

	"github.com/frantjc/sequence/pb/types"
)

type LogStreamServer interface {
	Send(*types.Log) error
}

func NewLogOutStreamWriter(stream LogStreamServer) io.Writer {
	return &logStreamWriter{sync.Mutex{}, stream, 0}
}

func NewLogErrStreamWriter(stream LogStreamServer) io.Writer {
	return &logStreamWriter{sync.Mutex{}, stream, 1}
}

func NewLogStreamMultiplexWriter(stream LogStreamServer) (io.Writer, io.Writer) {
	return NewLogOutStreamWriter(stream), NewLogErrStreamWriter(stream)
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
