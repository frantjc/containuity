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
	return &logOutStreamWriter{sync.Mutex{}, s}
}

func NewLogErrStreamWriter(s LogStreamServer) io.Writer {
	return &logErrStreamWriter{sync.Mutex{}, s}
}

func NewLogStreamMultiplexWriter(s LogStreamServer) (io.Writer, io.Writer) {
	return NewLogOutStreamWriter(s), NewLogErrStreamWriter(s)
}

type logOutStreamWriter struct {
	sync.Mutex
	s LogStreamServer
}

func (w *logOutStreamWriter) Write(p []byte) (int, error) {
	if p == nil {
		return len(p), nil
	}

	w.Lock()
	defer w.Unlock()
	err := w.s.Send(&types.Log{
		Out: p,
	})
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

type logErrStreamWriter struct {
	sync.Mutex
	s LogStreamServer
}

func (w *logErrStreamWriter) Write(p []byte) (int, error) {
	if p == nil {
		return len(p), nil
	}

	w.Lock()
	defer w.Unlock()
	err := w.s.Send(&types.Log{
		Err: p,
	})
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

var (
	_ io.Writer = &logOutStreamWriter{}
	_ io.Writer = &logErrStreamWriter{}
)
