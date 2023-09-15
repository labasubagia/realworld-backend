package logger

import (
	"sort"

	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

const defaultType = TypeZeroLog

var fnNewMap = map[string]func(util.Config) port.Logger{
	TypeZeroLog: NewZeroLogLogger,
	TypeZap:     NewZapLogger,
}

func Keys() (keys []string) {
	for key := range fnNewMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

func NewLogger(config util.Config) port.Logger {
	new, ok := fnNewMap[config.LogType]
	if ok {
		return new(config)
	}
	return fnNewMap[defaultType](config)
}
