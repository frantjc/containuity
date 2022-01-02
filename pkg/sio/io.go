package sio

import (
	"io"
	"os"
)

type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type IOOpt func(s *Streams)

func WithStdio(s *Streams) {
	s.Stdin = os.Stdin
	s.Stdout = os.Stdout
	s.Stderr = os.Stderr
}

func WithStreams(stdin io.Reader, stdout, stderr io.Writer) IOOpt {
	return func(s *Streams) {
		s.Stdin = stdin
		s.Stdout = stdout
		s.Stderr = stderr
	}
}

func New(opts ...IOOpt) *Streams {
	s := &Streams{}
	for _, opt := range opts {
		opt(s)
	}

	return s
}
