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
	log   zerolog.Logger
	event *zerolog.Event
}

func NewZeroLogLogger(config util.Config) port.Logger {
	logger := zerolog.New(os.Stderr)
	if !config.IsProduction() {
		output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}
	return &zeroLogLogger{
		log:   logger,
		event: logger.Info(),
	}
}

// Log Level

func (l *zeroLogLogger) Info() port.Logger {
	l.event = l.log.Info()
	return l
}

func (l *zeroLogLogger) Error() port.Logger {
	l.event = l.log.Error()
	return l
}

func (l *zeroLogLogger) Fatal() port.Logger {
	l.event = l.log.Fatal()
	return l
}

// Set Attributes

func (l *zeroLogLogger) Err(err error) port.Logger {
	l.event = l.event.Err(err)
	return l
}

func (l *zeroLogLogger) Field(key string, value any) port.Logger {
	l.event = l.event.Any(key, value)
	return l
}

// Send

func (l *zeroLogLogger) Msg(v ...any) {
	l.event.Msg(fmt.Sprint(v...))
}

func (l *zeroLogLogger) Msgf(format string, v ...any) {
	l.event.Msgf(format, v...)
}
