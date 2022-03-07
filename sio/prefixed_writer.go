package sio

import (
	"io"
)

const (
	newline = '\n'
)

type prefixedWriter struct {
	p []byte
	w io.Writer
}

var _ io.Writer = &prefixedWriter{}

func (w *prefixedWriter) Write(p []byte) (c int, err error) {
	var (
		a int
		j = 0
	)
	for i, b := range p {
		if b == newline {
			a, err = w.w.Write(append(w.p, p[j:i+1]...))
			c += a
			if err != nil {
				return
			}
			j = i + 1
		}
	}

	if s := p[j:]; len(s) > 0 {
		a, err = w.w.Write(append(append(w.p, s...), newline))
		c += a
		if err != nil {
			return
		}
	}

	return a, nil
}

func NewPrefixedWriter(prefix string, w io.Writer) io.Writer {
	p := []byte(prefix)
	return &prefixedWriter{p, w}
}
