package sio

import "io"

type noOpWriter struct{}

var _ io.Writer = &noOpWriter{}

func (w *noOpWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func NewNoOpWriter() io.Writer {
	return &noOpWriter{}
}
