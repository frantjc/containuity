package log

import (
	"fmt"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
)

func fmtNoOp(i interface{}) string {
	return ""
}

func fmtMsg(i interface{}) string {
	return fmt.Sprintf("%s", i)
}

func fmtField(i interface{}) string {
	s := fmt.Sprintf("%s", i)
	if strings.Contains(s, " ") {
		s = fmt.Sprintf("\"%s\"", s)
	}
	return s
}

func init() {
	logger = zerolog.New(
		zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) {
				w.FormatTimestamp = fmtNoOp
				w.FormatLevel = fmtNoOp
				w.FormatMessage = fmtMsg
				w.FormatFieldName = fmtField
				w.FormatFieldValue = fmtField
				w.Out = colorable.NewColorableStdout()
			},
		),
	).Level(zerolog.DebugLevel)
}
