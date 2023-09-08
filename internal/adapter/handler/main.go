package handler

import (
	"github.com/labasubagia/realworld-backend/internal/adapter/handler/rest"
	"github.com/labasubagia/realworld-backend/internal/core/port"
)

func NewServer(service port.Service) port.Server {
	return rest.NewServer(service)
}
