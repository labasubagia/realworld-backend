package handler

import (
	"github.com/labasubagia/realworld-backend/internal/adapter/handler/restful"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

func NewServer(config util.Config, service port.Service, logger port.Logger) port.Server {
	return restful.NewServer(config, service, logger)
}
