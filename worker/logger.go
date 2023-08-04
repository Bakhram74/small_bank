package worker

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}
func (logger *Logger) print(level zerolog.Level, arg ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(arg...))
}

func (logger *Logger) Info(args ...interface{}) {
	logger.print(zerolog.DebugLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.print(zerolog.WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.print(zerolog.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.print(zerolog.FatalLevel, args...)
}
func (logger *Logger) Debug(args ...interface{}) {
	logger.print(zerolog.DebugLevel, args...)
}
