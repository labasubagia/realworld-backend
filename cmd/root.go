package cmd

import (
	"log"

	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/spf13/cobra"
)

var config util.Config

var rootCmd = &cobra.Command{
	Use:   "realworld",
	Short: "Realworld backend app",
	Long:  "Realworld is an app about article similar medium.com and dev.to",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	var err error

	config, err = util.LoadConfig(".env")
	if err != nil {
		log.Fatal("failed to load env config", err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
