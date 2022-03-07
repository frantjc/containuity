package actions

import "io"

type commandWriter struct {
	callback func(*Command) []byte
	w        io.Writer
}

var _ io.Writer = &commandWriter{}

func (w *commandWriter) Write(p []byte) (int, error) {
	if c, err := ParseCommand(p); err == nil {
		if b := w.callback(c); len(b) == 0 {
			return len(p), nil
		} else {
			_, err = w.w.Write(b)
			return len(p), err
		}
	}
	return w.w.Write(p)
}

func NewCommandWriter(callback func(*Command) []byte, w io.Writer) io.Writer {
	return &commandWriter{callback, w}
}
