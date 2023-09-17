package cmd

import (
	"fmt"
	"strings"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler"
	"github.com/labasubagia/realworld-backend/internal/adapter/logger"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository"
	"github.com/labasubagia/realworld-backend/internal/core/service"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/spf13/cobra"
)

func init() {
	dbTypeStr := strings.Join(repository.Keys(), ", ")
	logTypeStr := strings.Join(logger.Keys(), ", ")
	serverTypeStr := strings.Join(handler.Keys(), ", ")

	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Bool("prod", config.IsProduction(), "use for production")
	serverCmd.Flags().StringVarP(&config.ServerType, "server", "s", config.ServerType, fmt.Sprintf("server type in (%s)", serverTypeStr))
	serverCmd.Flags().IntVarP(&config.ServerPort, "port", "p", config.ServerPort, "server port number")
	serverCmd.Flags().StringVarP(&config.DBType, "database", "d", config.DBType, fmt.Sprintf("database type in (%s)", dbTypeStr))
	serverCmd.Flags().StringVarP(&config.LogType, "log", "l", config.LogType, fmt.Sprintf("log type in (%s)", logTypeStr))
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run http server",
	Long:  "Run gin server restful API",
	Run: func(cmd *cobra.Command, args []string) {

		// is_prod
		isProduction, err := cmd.Flags().GetBool("prod")
		if err == nil && isProduction {
			config.Environment = util.EnvProduction
		}
		logger := logger.NewLogger(config)
		logger.Info().Msgf("use logger %s", config.LogType)

		// repository
		repo, err := repository.NewRepository(config, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to load repository")
		}
		logger.Info().Msgf("use repository %s", config.DBType)

		service, err := service.NewService(config, repo, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to load service")
		}

		logger.Info().Msgf("%s server listen to port %d", config.ServerType, config.ServerPort)
		server := handler.NewServer(config, service, logger)
		if server.Start(); err != nil {
			logger.Fatal().Err(err).Msg("failed to load service")
		}
	},
}
