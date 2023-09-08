package cmd

import (
	"log"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository"
	"github.com/labasubagia/realworld-backend/internal/core/service"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/spf13/cobra"
)

func startCmd(config util.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "This start restful API server",
		Run: func(cmd *cobra.Command, args []string) {
			repo, err := repository.NewRepository(config)
			if err != nil {
				log.Fatal("failed to init repository", err)
			}
			svc := service.NewService(repo)
			server := handler.NewServer(config, svc)
			if err := server.Start(); err != nil {
				log.Fatal("failed to start server", err)
			}
		},
	}
}
