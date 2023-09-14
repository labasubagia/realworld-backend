package port

type Logger interface {
	// Level
	// ? just limit to 2
	Info() Logger
	Error() Logger
	Fatal() Logger
	Err(error) Logger

	// set attributes
	Field(string, any) Logger

	// send
	Msgf(string, ...any)
	Msg(...any)
}
