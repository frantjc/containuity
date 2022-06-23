package rpcio

import (
	"errors"
	"io"
)

func NewServerStreamWriter[T any](serverStream ServerStream[T], convert func([]byte) *T) io.Writer {
	return &ServerStreamWriter[T]{serverStream, convert}
}

type ServerStream[T any] interface {
	Send(*T) error
	CloseAndReceive() (any, error)
}

type ServerStreamWriter[T any] struct {
	ServerStream ServerStream[T]
	Convert      func([]byte) *T
}

var _ io.Writer = &ServerStreamWriter[any]{}

func (w *ServerStreamWriter[T]) Write(p []byte) (int, error) {
	err := w.ServerStream.Send(w.Convert(p))
	switch {
	case errors.Is(err, io.EOF):
		_, err = w.ServerStream.CloseAndReceive()
		return 0, err
	case err != nil:
		return 0, err
	}

	return len(p), nil
}
