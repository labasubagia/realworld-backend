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
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP(
		"port",
		"p",
		config.HTTPServerPort,
		"server port number",
	)
	serverCmd.Flags().StringP(
		"database",
		"d",
		repository.DefaultRepoKey,
		fmt.Sprintf("database type in (%s)", strings.Join(repository.Keys(), ", ")),
	)
	serverCmd.Flags().StringP(
		"log",
		"l",
		logger.DefaultKey,
		fmt.Sprintf("log type in (%s)", strings.Join(logger.Keys(), ",")),
	)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run http server",
	Long:  "Run gin server restful API",
	Run: func(cmd *cobra.Command, args []string) {

		logType := cmd.Flag("log").Value.String()
		log := logger.FnNewMap[logType](config)

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get flag port")
		}
		config.HTTPServerPort = port

		dbType, err := cmd.Flags().GetString("database")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get flag database")
		}
		dbType = strings.ToLower(dbType)

		newRepo, exist := repository.FnNewMap[dbType]
		if !exist {
			fmt.Println()
			log.Fatal().Err(err).Msg("failed to get flag port")
		}

		repo, err := newRepo(config, log)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load repository")
		}

		service, err := service.NewService(config, repo, log)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load service")
		}

		server := restful.NewServer(config, service, log)
		if server.Start(); err != nil {
			log.Fatal().Err(err).Msg("failed to load service")
		}
	},
}
