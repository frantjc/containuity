package grpcio

import "errors"

var (
	ErrStreamClosed = errors.New("stream closed")
)

func ErrIsStreamClosed(err error) bool {
	return errors.Is(err, ErrStreamClosed)
}
