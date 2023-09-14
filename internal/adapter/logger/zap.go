package logger

import (
	"fmt"

	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	config util.Config
	level  zapcore.Level
	fields map[string]any
}

func NewZapLogger(config util.Config) port.Logger {
	return &zapLogger{
		config: config,
		fields: map[string]any{},
		level:  zap.DebugLevel,
	}
}

func (l *zapLogger) Field(key string, value any) port.Logger {
	l.fields[key] = value
	return l
}

func (l *zapLogger) Logger() port.Logger {
	return l
}

func (l *zapLogger) Info() port.LogEvent {
	l.level = zap.InfoLevel
	return newZapEvent(l)
}

func (l *zapLogger) Error() port.LogEvent {
	l.level = zapcore.ErrorLevel
	return newZapEvent(l)
}

func (l *zapLogger) Fatal() port.LogEvent {
	l.level = zap.PanicLevel
	return newZapEvent(l)
}

type zapEvent struct {
	opt    *zapLogger
	fields map[string]any
}

func newZapEvent(opt *zapLogger) port.LogEvent {
	return &zapEvent{
		opt:    opt,
		fields: map[string]any{},
	}
}

func (e *zapEvent) Err(err error) port.LogEvent {
	if err != nil {
		return e
	}
	e.fields["error"] = err.Error()
	return e
}

func (e *zapEvent) Field(key string, value any) port.LogEvent {
	e.fields[key] = value
	return e
}

func (e *zapEvent) Msg(v ...any) {
	msg := fmt.Sprint(v...)
	e.send(msg)
}

func (e *zapEvent) Msgf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	e.send(msg)
}

func (e *zapEvent) send(msg string) {
	config := zap.NewProductionConfig()
	if !e.opt.config.IsProduction() {
		config = zap.NewDevelopmentConfig()
	}

	for k, v := range e.opt.fields {
		e.fields[k] = v
	}
	e.opt.fields = map[string]any{}
	config.InitialFields = e.fields

	logger, _ := config.Build()
	defer func() {
		logger.Sync()
		e.fields = map[string]any{}
	}()
	stdLogger, _ := zap.NewStdLogAt(logger, e.opt.level)
	stdLogger.Println(msg)
}
