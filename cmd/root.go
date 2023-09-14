package cmd

import (
	"fmt"
	"os"

	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/spf13/cobra"
)

var config util.Config

var rootCmd = &cobra.Command{
	Use:   "realworld",
	Short: "Realworld backend app",
	Long:  "Realworld is an app about article similar medium.com and dev.to",
}

func init() {
	var err error

	config, err = util.LoadConfig(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load env config: %s", err)
		os.Exit(1)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
