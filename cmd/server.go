package cmd

import (
	"fmt"
	"strings"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler/restful"
	"github.com/labasubagia/realworld-backend/internal/adapter/logger"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository"
	"github.com/labasubagia/realworld-backend/internal/core/service"
	"github.com/spf13/cobra"
)

func init() {
	dbOpts := []string{}
	for option := range repository.RepoFnMap {
		dbOpts = append(dbOpts, option)
	}
	dbOptStr := strings.Join(dbOpts, ",")

	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringP("database", "d", repository.DefaultRepoKey, fmt.Sprintf("select database in (%s)", dbOptStr))
	serverCmd.Flags().IntP("port", "p", config.HTTPServerPort, "server port")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run http server",
	Long:  "Run gin server restful API",
	Run: func(cmd *cobra.Command, args []string) {

		logger := logger.NewLogger(config)

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to get flag port")
		}
		config.HTTPServerPort = port

		dbType, err := cmd.Flags().GetString("database")
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to get flag database")
		}
		dbType = strings.ToLower(dbType)

		newRepo, exist := repository.RepoFnMap[dbType]
		if !exist {
			fmt.Println()
			logger.Fatal().Err(err).Msg("failed to get flag port")
		}

		repo, err := newRepo(config, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to load repository")
		}

		service, err := service.NewService(config, repo, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to load service")
		}

		server := restful.NewServer(config, service, logger)
		if server.Start(); err != nil {
			logger.Fatal().Err(err).Msg("failed to load service")
		}
	},
}
