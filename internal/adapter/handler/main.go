package handler

import (
	"sort"

	grpc_api "github.com/labasubagia/realworld-backend/internal/adapter/handler/grpc/api"
	"github.com/labasubagia/realworld-backend/internal/adapter/handler/restful"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

const defaultType = restful.TypeRestful

var fnNewMap = map[string]func(util.Config, port.Service, port.Logger) port.Server{
	restful.TypeRestful: restful.NewServer,
	grpc_api.TypeGrpc:   grpc_api.NewServer,
}

func Keys() (keys []string) {
	for key := range fnNewMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

func NewServer(config util.Config, service port.Service, logger port.Logger) port.Server {
	new, ok := fnNewMap[config.ServerType]
	if ok {
		return new(config, service, logger)
	}
	return fnNewMap[defaultType](config, service, logger)
}
