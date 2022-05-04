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
	if p == nil || len(p) == 0 {
		return 0, nil
	}

	lines := strings.Split(string(p), "\n")
	for i, line := range lines {
		a := []byte{}
		if i < len(lines)-1 {
			a = []byte{'\n'}
		}

		if len(line) == 0 {
		} else if c, err := ParseStringCommand(line); err == nil {
			if b := w.callback(c); len(b) != 0 {
				if _, err = w.w.Write(append(b, a...)); err != nil {
					return len(p), err
				}
			}
		} else {
			if _, err := w.w.Write(append([]byte(line), a...)); err != nil {
				return len(p), err
			}
		}
	}

	return len(p), nil
}

func NewCommandWriter(callback func(*Command) []byte, w io.Writer) io.Writer {
	return &commandWriter{callback, w}
}
