package log

import "io"

type prefixedWriter struct {
	prefix string
	callback func (s string, v... interface{})
}

var _ io.Writer = &prefixedWriter{}

func (w *prefixedWriter) Write(p []byte) (int, error) {
	j := 0
	for i, b := range p {
		if b == '\n' {
			w.callback("%s%s", w.prefix, p[j:i])
			j = i+1
		}
	}

	if s := p[j:]; len(s) > 0 {
		w.callback("%s%s", w.prefix, p[j:])
	}

	return len(p), nil
}

func NewPrefixedDebugWriter(p string) io.Writer {
	return &prefixedWriter{
		prefix: p,
		callback: Debugf,
	}
}

func NewPrefixedInfoWriter(p string) io.Writer {
	return &prefixedWriter{
		prefix: p,
		callback: Infof,
	}
}
