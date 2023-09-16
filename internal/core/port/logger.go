package port

import (
	"context"
)

// set context value with this key
const SubLoggerCtxKey = "sub_logger"

type Logger interface {
	NewInstance() Logger

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
