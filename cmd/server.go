package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler/restful"
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

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatal("failed get port flag", err)
		}
		config.HTTPServerPort = port

		dbType, err := cmd.Flags().GetString("database")
		if err != nil {
			log.Fatal("failed get database flag", err)
		}
		dbType = strings.ToLower(dbType)

		newRepo, exist := repository.RepoFnMap[dbType]
		if !exist {
			log.Fatal("invalid database", dbType)
		}

		repo, err := newRepo(config)
		if err != nil {
			log.Fatal("failed to load repository", err)
		}
		service, err := service.NewService(config, repo)
		if err != nil {
			log.Fatal("failed to load service", err)
		}
		server := restful.NewServer(config, service)

		if server.Start(); err != nil {
			log.Fatal("failed to load server", err)
		}
	},
}
