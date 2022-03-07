package log

import "io"

func Writer() io.Writer {
	return logger
}
