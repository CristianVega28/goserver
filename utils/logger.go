package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type (
	LoggerI interface {
		Fatal(message string)
		Msg(message string)
		Create() Logger
	}

	Logger struct {
		log zerolog.Logger
	}
)

func (l *Logger) Create() Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return Logger{
		log: zerolog.New(output).With().Timestamp().Logger(),
	}
}

func (l *Logger) Fatal(message string) {
	l.log.Fatal().Msg(message)
}

func (l *Logger) Msg(message string) {
	l.log.Debug().Msg(message)
}
