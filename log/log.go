package log

import (
	"io"
	"os"
)

var (
	log = New(os.Stdout)
)

// for some reason, returning log here
// instead of log.l breaks everything
func Writer() io.Writer {
	return log
}

func SetVerbose(v bool) {
	log.SetVerbose(v)
}

func Debug(s string) {
	log.Debug(s)
}

func Debugf(s string, v ...interface{}) {
	log.Debugf(s, v...)
}

func Info(s string) {
	log.Info(s)
}

func Infof(s string, v ...interface{}) {
	log.Infof(s, v...)
}
