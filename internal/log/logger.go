package log

import (
	"fmt"
	"io"
	"strings"

	"github.com/frantjc/go-js"
	"github.com/rs/zerolog"
)

type Logger interface {
	io.Writer
	SetVerbose(bool) Logger
	Debug(string)
	Debugf(string, ...interface{})
	Info(string)
	Infof(string, ...interface{})
}

type logger struct {
	zerolog zerolog.Logger
}

func (l *logger) Write(p []byte) (int, error) {
	return l.zerolog.Write(p)
}

func (l *logger) SetVerbose(v bool) Logger {
	l.zerolog = l.zerolog.Level(js.Ternary(v, zerolog.DebugLevel, zerolog.InfoLevel))
	return l
}

func (l *logger) Debug(s string) {
	l.zerolog.Debug().Msg(s)
}

func (l *logger) Debugf(s string, v ...interface{}) {
	l.zerolog.Debug().Msgf(s, v...)
}

func (l *logger) Info(s string) {
	l.zerolog.Info().Msg(s)
}

func (l *logger) Infof(s string, v ...interface{}) {
	l.zerolog.Info().Msgf(s, v...)
}

func New(o io.Writer) Logger {
	return &logger{
		zerolog: zerolog.New(
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
	return fmt.Sprint(i)
}

func fmtField(i interface{}) string {
	s := fmtMsg(i)
	if strings.Contains(s, " ") {
		s = fmt.Sprintf("\"%s\"", s)
	}
	return s
}
