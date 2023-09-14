package port

import (
	"context"
)

const SubLoggerCtxKey = "logger_key"

type Logger interface {
	Field(string, any) Logger
	Logger() Logger

	Info() LogEvent
	Error() LogEvent
	Fatal() LogEvent
}

type LogEvent interface {
	Field(string, any) LogEvent
	Err(error) LogEvent

	Msgf(string, ...any)
	Msg(...any)
}

func GetCtxSubLogger(ctx context.Context, defaultLogger Logger) Logger {
	if subLogger, ok := ctx.Value(SubLoggerCtxKey).(Logger); ok {
		return subLogger
	}
	return defaultLogger
}
