package log

import (
	"io"

	"github.com/mattn/go-colorable"
)

var (
	log = New(colorable.NewColorableStdout())
)

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
