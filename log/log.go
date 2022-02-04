package log

import "github.com/rs/zerolog"

var (
	logger zerolog.Logger
)

func SetVerbose(v bool) {
	if v {
		logger = logger.Level(zerolog.DebugLevel)
	} else {
		logger = logger.Level(zerolog.InfoLevel)
	}
}

func Debug(s string) {
	logger.Debug().Msg(s)
}

func Debugf(s string, v... interface{}) {
	logger.Debug().Msgf(s, v...)
}

func Info(s string) {
	logger.Info().Msg(s)
}

func Infof(s string, v... interface{}) {
	logger.Info().Msgf(s, v...)
}
