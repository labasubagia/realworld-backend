package handler

import (
	"github.com/labasubagia/realworld-backend/internal/adapter/handler/rest"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

func NewServer(config util.Config, service port.Service) port.Server {
	return rest.NewServer(config, service)
}
