package cmd

import (
	"fmt"

	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/spf13/cobra"
)

type Command struct {
	config  util.Config
	rootCmd *cobra.Command
}

func (c *Command) AddCommand(fnCommands ...func(config util.Config) *cobra.Command) {
	commands := []*cobra.Command{}
	for _, fn := range fnCommands {
		commands = append(commands, fn(c.config))
	}
	c.rootCmd.AddCommand(commands...)
}

func (c *Command) Execute() error {
	return c.rootCmd.Execute()
}

func NewCommand(config util.Config) *Command {
	command := &Command{
		config: config,
		rootCmd: &cobra.Command{
			Use:   "app",
			Short: "App about article",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("see help")
			},
		},
	}
	command.AddCommand(startCmd)
	return command
}
