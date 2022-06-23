package rpcio

import (
	"bytes"
	"io"
)

func NewClientStreamReadCloser[T any](clientStream ClientStream[T], convert func(*T) []byte) io.ReadCloser {
	return &ClientStreamReadCloser[T]{clientStream, new(bytes.Buffer), convert}
}

type ClientStream[T any] interface {
	Msg() *T
	Err() error
	Receive() bool
	Close() error
}

type Buffer interface {
	io.ReadWriter
	Len() int
}

type ClientStreamReadCloser[T any] struct {
	ClientStream ClientStream[T]
	Buffer       Buffer
	Convert      func(*T) []byte
}

var _ io.ReadCloser = &ClientStreamReadCloser[any]{}

func (r *ClientStreamReadCloser[T]) Read(p []byte) (int, error) {
	var (
		pLen = len(p)
	)
	for r.Buffer.Len() < pLen && r.ClientStream.Receive() {
		var (
			msg = r.ClientStream.Msg()
			err = r.ClientStream.Err()
		)
		if msg == nil || err != nil {
			return 0, io.ErrClosedPipe
		}

		if _, err = r.Buffer.Write(r.Convert(msg)); err != nil {
			return 0, io.ErrClosedPipe
		}
	}

	return r.Buffer.Read(p)
}

func (r *ClientStreamReadCloser[T]) Close() error {
	return r.ClientStream.Close()
}
