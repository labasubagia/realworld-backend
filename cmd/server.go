package cmd

import (
	"fmt"
	"strings"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler/restful"
	"github.com/labasubagia/realworld-backend/internal/adapter/logger"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository"
	"github.com/labasubagia/realworld-backend/internal/core/service"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/spf13/cobra"
)

func init() {
	dbTypeStr := strings.Join(repository.Keys(), ", ")
	logTypeStr := strings.Join(logger.Keys(), ", ")

	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Bool("prod", false, "use for production")
	serverCmd.Flags().IntVarP(&config.HTTPServerPort, "port", "p", config.HTTPServerPort, "server port number")
	serverCmd.Flags().StringVarP(&config.DBType, "database", "d", repository.DefaultType, fmt.Sprintf("database type in (%s)", dbTypeStr))
	serverCmd.Flags().StringVarP(&config.LogType, "log", "l", logger.DefaultType, fmt.Sprintf("log type in (%s)", logTypeStr))
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

		logger.Info().Msgf("listen to port %d", config.HTTPServerPort)
		server := restful.NewServer(config, service, logger)
		if server.Start(); err != nil {
			logger.Fatal().Err(err).Msg("failed to load service")
		}
	},
}
