package logger

import (
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

func NewLogger(config util.Config) port.Logger {
	return NewZeroLogLogger(config)
}
