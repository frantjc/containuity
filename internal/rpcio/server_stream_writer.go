package rpcio

import "io"

func NewServerStreamWriter[T any](serverStream ServerStream[T], convert func([]byte) *T) io.Writer {
	return &ServerStreamWriter[T]{serverStream, convert}
}

type ServerStream[T any] interface {
	Send(*T) error
}

type ServerStreamWriter[T any] struct {
	ServerStream ServerStream[T]
	Convert      func([]byte) *T
}

func (w *ServerStreamWriter[T]) Write(p []byte) (int, error) {
	if err := w.ServerStream.Send(w.Convert(p)); err != nil {
		return 0, err
	}

	return len(p), nil
}
