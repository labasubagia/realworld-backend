package logger

import (
	"sort"

	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

type FnNew func(util.Config) port.Logger

const DefaultType = TypeZeroLog

var FnNewMap = map[string]FnNew{
	TypeZeroLog: NewZeroLogLogger,
	TypeZap:     NewZapLogger,
}

func Keys() (keys []string) {
	for key := range FnNewMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

func NewLogger(config util.Config) port.Logger {
	return FnNewMap[DefaultType](config)
}
