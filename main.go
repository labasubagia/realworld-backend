package main

import (
	"log"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler"
	repository "github.com/labasubagia/realworld-backend/internal/adapter/repository"
	"github.com/labasubagia/realworld-backend/internal/core/service"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

func main() {
	config, err := util.LoadConfig(".env")
	if err != nil {
		log.Fatal("failed to load config", err)
	}
	repo, err := repository.NewRepository(config)
	if err != nil {
		log.Fatal("failed to init repository", err)
	}
	svc := service.NewService(repo)
	server := handler.NewServer(svc)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
