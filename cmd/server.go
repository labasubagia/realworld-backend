package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler/restful"
	"github.com/labasubagia/realworld-backend/internal/adapter/logger"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository"
	"github.com/labasubagia/realworld-backend/internal/core/service"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/spf13/cobra"
)

var (
	dbTypeStr  = strings.Join(repository.Keys(), ", ")
	logTypeStr = strings.Join(logger.Keys(), ", ")
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Bool("prod", false, "use for production")
	serverCmd.Flags().IntP("port", "p", config.HTTPServerPort, "server port number")
	serverCmd.Flags().StringP("database", "d", repository.DefaultType, fmt.Sprintf("database type in (%s)", dbTypeStr))
	serverCmd.Flags().StringP("log", "l", logger.DefaultType, fmt.Sprintf("log type in (%s)", logTypeStr))
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

		// port
		port, err := cmd.Flags().GetInt("port")
		if err == nil {
			config.HTTPServerPort = port
		}

		// logger
		logType, err := cmd.Flags().GetString("log")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed get log type flag: %s", err)
			os.Exit(1)
		}
		newLogger, ok := logger.FnNewMap[logType]
		if !ok {
			logType = logger.DefaultType
			newLogger = logger.NewLogger
		}
		logger := newLogger(config)
		logger.Info().Msgf("use logger %s", logType)

		// db
		dbType, err := cmd.Flags().GetString("database")
		if err != nil {
			logger.Fatal().Err(err).Msg("failed get log type flag")
		}
		dbType = strings.ToLower(dbType)
		newRepo, ok := repository.FnNewMap[dbType]
		if !ok {
			dbType = repository.DefaultType
			newRepo = repository.NewRepository
		}
		logger.Info().Msgf("use database %s", dbType)

		repo, err := newRepo(config, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to load repository")
		}

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
