package logger

import (
	"sort"

	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

const DefaultKey = "default"

type FnNew func(util.Config) port.Logger

var FnNewMap = map[string]FnNew{
	DefaultKey: NewZeroLogLogger,
	"zerolog":  NewZeroLogLogger,
	"zap":      NewZapLogger,
}

func LogKeys() (keys []string) {
	for key := range FnNewMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

func NewLogger(config util.Config) port.Logger {
	return FnNewMap[DefaultKey](config)
}
