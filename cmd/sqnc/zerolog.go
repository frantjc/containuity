package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/frantjc/sequence"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	o := zerolog.ConsoleWriter{
		Out:        sequence.Stdout,
		TimeFormat: time.RFC3339Nano,
	}

	o.FormatTimestamp = func(i interface{}) string {
		return fmt.Sprintf("time=\"%s\"", i)
	}

	o.FormatLevel = func(i interface{}) string {
		return fmt.Sprintf("level=%s", i)
	}

	o.FormatMessage = func(i interface{}) string {
		s := fmt.Sprintf("%s", i)
		if strings.Contains(s, " ") {
			s = fmt.Sprintf("msg=\"%s\"", s)
		}
		return s
	}

	o.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}

	o.FormatFieldValue = func(i interface{}) string {
		s := fmt.Sprintf("%s", i)
		if strings.Contains(s, " ") {
			s = fmt.Sprintf("\"%s\"", s)
		}
		return s
	}

	log.Logger = zerolog.New(o).With().Timestamp().Logger()
}
