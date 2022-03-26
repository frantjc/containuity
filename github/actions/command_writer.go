package actions

import (
	"io"
	"strings"
)

type commandWriter struct {
	callback func(*Command) []byte
	w        io.Writer
}

var _ io.Writer = &commandWriter{}

func (w *commandWriter) Write(p []byte) (int, error) {
	for _, line := range strings.Split(string(p), "\n") {
		if c, err := ParseStringCommand(line); err == nil {
			if b := w.callback(c); len(b) != 0 {
				if _, err = w.w.Write(b); err != nil {
					return len(p), err
				}
			}
		} else {
			if _, err := w.w.Write([]byte(line)); err != nil {
				return len(p), err
			}
		}
	}
	return len(p), nil
}

func NewCommandWriter(callback func(*Command) []byte, w io.Writer) io.Writer {
	return &commandWriter{callback, w}
}
