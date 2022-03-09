package log

import (
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog"
)

type Logger interface {
	io.Writer
	SetVerbose(bool)
	Debug(string)
	Debugf(string, ...interface{})
	Info(string)
	Infof(string, ...interface{})
}

type logger struct {
	l zerolog.Logger
}

var _ Logger = &logger{}

func (l *logger) Write(p []byte) (int, error) {
	return l.l.Write(p)
}

func (l *logger) SetVerbose(v bool) {
	if v {
		l.l = l.l.Level(zerolog.DebugLevel)
	} else {
		l.l = l.l.Level(zerolog.InfoLevel)
	}
}

func (l *logger) Debug(s string) {
	l.l.Debug().Msg(s)
}

func (l *logger) Debugf(s string, v ...interface{}) {
	l.l.Debug().Msgf(s, v...)
}

func (l *logger) Info(s string) {
	l.l.Info().Msg(s)
}

func (l *logger) Infof(s string, v ...interface{}) {
	l.l.Info().Msgf(s, v...)
}

func New(o io.Writer) Logger {
	return &logger{
		l: zerolog.New(
			zerolog.NewConsoleWriter(
				func(w *zerolog.ConsoleWriter) {
					w.FormatTimestamp = fmtNoOp
					w.FormatLevel = fmtNoOp
					w.FormatMessage = fmtMsg
					w.FormatFieldName = fmtField
					w.FormatFieldValue = fmtField
					w.Out = o
				},
			),
		).Level(zerolog.InfoLevel),
	}
}

func fmtNoOp(i interface{}) string {
	return ""
}

func fmtMsg(i interface{}) string {
	return fmt.Sprintf("%s", i)
}

func fmtField(i interface{}) string {
	s := fmtMsg(i)
	if strings.Contains(s, " ") {
		s = fmt.Sprintf("\"%s\"", s)
	}
	return s
}
