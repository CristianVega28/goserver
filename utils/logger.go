package utils

import (
	"fmt"
	"os"
	"runtime"
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
	_, file, line, _ := runtime.Caller(1)
	l.log.Fatal().Msg(fmt.Sprintf("Fatal: File %s - %d | Msg: %s", file, line, message))
}

func (l *Logger) Msg(message string) {
	l.log.Debug().Msg(message)
}

func (l *Logger) Everyone(message string, maps any) {
	l.log.Info().Fields(maps).Msg(message)
}

// Structs and slice logs
func (l *Logger) Structs(message string, structs any) {
	l.log.Info().Interface("obj", structs).Msg(message)
}

func (l *Logger) Slice(message string, array any) {
	l.log.Info().Interface(message, array)
}

var Log Logger

func InitLogger() {
	Log = (&Logger{}).Create()

}
