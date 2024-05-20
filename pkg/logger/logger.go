package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

type Logger interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type logger struct {
	logger *zerolog.Logger
}

func New(level string) Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)
	skipFrameCount := 3
	var z zerolog.Logger
	if l == zerolog.DebugLevel {

		z = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
			Logger()
	} else {
		z = zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
			Logger()
	}

	return &logger{
		logger: &z,
	}
}

func (l *logger) Debug(message any, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Debug(), mf, args...)
}

func (l *logger) Info(message string, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Info(), mf, args...)
}

func (l *logger) Warn(message string, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Warn(), mf, args...)
}

func (l *logger) Error(message interface{}, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Error(), mf, args...)
}

func (l *logger) Fatal(message interface{}, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Fatal(), mf, args...)
	os.Exit(1)
}
func (l *logger) formatMessage(message any) string {
	switch t := message.(type) {
	case error:
		return t.Error()
	case string:
		return t
	default:
		return fmt.Sprintf("Unknown type %v", message)
	}
}
func (l *logger) log(e *zerolog.Event, m string, args ...any) {
	if len(args) == 0 {
		e.Msg(m)
	} else {
		e.Msgf(m, args...)
	}
}
