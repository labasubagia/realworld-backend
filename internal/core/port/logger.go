package port

const LoggerCtxKey = "logger_key"

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
