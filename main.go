package main

import (
	"log"

	"github.com/labasubagia/realworld-backend/api"
	"github.com/labasubagia/realworld-backend/repository"
	"github.com/labasubagia/realworld-backend/service"
	"github.com/labasubagia/realworld-backend/util"
)

func main() {
	config, err := util.LoadConfig(".env")
	if err != nil {
		log.Fatal("failed to load config", err)
	}
	repo, err := repository.NewSQLRepository(config)
	if err != nil {
		log.Fatal("failed to init repository", err)
	}
	svc := service.NewService(repo)
	server := api.NewServer(svc)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
