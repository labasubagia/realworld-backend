package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/rs/zerolog"
)

type zeroLogLogger struct {
	logger zerolog.Logger
	fields map[string]any
}

func NewZeroLogLogger(config util.Config) port.Logger {
	logger := zerolog.New(os.Stderr)
	if !config.IsProduction() {
		output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}
	return &zeroLogLogger{
		logger: logger,
		fields: map[string]any{},
	}
}

func (l *zeroLogLogger) Field(key string, value any) port.Logger {
	l.fields[key] = value
	return l
}

func (l *zeroLogLogger) Logger() port.Logger {
	return l
}

func (l *zeroLogLogger) Info() port.LogEvent {
	return newZeroLogEvent(l.fields, l.logger.Info())
}

func (l *zeroLogLogger) Error() port.LogEvent {
	return newZeroLogEvent(l.fields, l.logger.Error())
}

func (l *zeroLogLogger) Fatal() port.LogEvent {
	return newZeroLogEvent(l.fields, l.logger.Fatal())
}

type zeroLogEvent struct {
	event *zerolog.Event
}

func newZeroLogEvent(initialFields map[string]any, event *zerolog.Event) port.LogEvent {
	event.Fields(initialFields)
	return &zeroLogEvent{event: event}
}

func (e *zeroLogEvent) Err(err error) port.LogEvent {
	e.event.Err(err)
	return e
}

func (e *zeroLogEvent) Field(key string, value any) port.LogEvent {
	e.event.Any(key, value)
	return e
}

func (e *zeroLogEvent) Msg(v ...any) {
	msg := fmt.Sprint(v...)
	e.event.Msg(msg)
}

func (e *zeroLogEvent) Msgf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	e.event.Msg(msg)
}
