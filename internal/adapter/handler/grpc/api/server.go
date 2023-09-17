package api

import (
	"fmt"
	"net"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler/grpc/pb"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const TypeGrpc = "grpc"

type Server struct {
	pb.UnimplementedRealWorldServer
	config  util.Config
	service port.Service
	logger  port.Logger
}

func NewServer(config util.Config, repository port.Service, logger port.Logger) port.Server {
	server := &Server{
		config:  config,
		service: repository,
		logger:  logger,
	}
	return server
}

func (server *Server) Start() error {
	logger := grpc.UnaryInterceptor(server.Logger)

	grpcServer := grpc.NewServer(logger)
	pb.RegisterRealWorldServer(grpcServer, server)
	reflection.Register(grpcServer)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", server.config.ServerPort))
	if err != nil {
		return err
	}
	if err := grpcServer.Serve(listen); err != nil {
		return err
	}
	return nil
}
