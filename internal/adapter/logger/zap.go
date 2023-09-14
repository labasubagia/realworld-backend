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

// Log Level

func (l *zapLogger) Info() port.Logger {
	l.level = zap.InfoLevel
	return l
}

func (l *zapLogger) Error() port.Logger {
	l.level = zapcore.ErrorLevel
	return l
}

func (l *zapLogger) Fatal() port.Logger {
	l.level = zap.PanicLevel
	return l
}

// Set Attributes

func (l *zapLogger) Field(key string, value any) port.Logger {
	l.fields[key] = value
	return l
}

func (l *zapLogger) Err(err error) port.Logger {
	if err != nil {
		return l
	}
	l.fields["error"] = err.Error()
	return l
}

// Send

func (l *zapLogger) Msg(v ...any) {
	msg := fmt.Sprintln(v...)
	l.send(msg)
}

func (l *zapLogger) Msgf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.send(msg)
}

func (l *zapLogger) send(msg string) {
	config := zap.NewProductionConfig()
	if !l.config.IsProduction() {
		config = zap.NewDevelopmentConfig()
	}
	config.InitialFields = l.fields
	logger, _ := config.Build()
	defer func() {
		logger.Sync()
		l.fields = map[string]any{}
	}()
	stdLogger, _ := zap.NewStdLogAt(logger, l.level)
	stdLogger.Println(msg)
}
