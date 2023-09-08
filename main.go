package main

import (
	"log"

	"github.com/labasubagia/realworld-backend/cmd"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

func main() {
	config, err := util.LoadConfig(".env")
	if err != nil {
		log.Fatal("failed to load config", err)
	}
	// run with go run main.go help
	command := cmd.NewCommand(config)
	if err := command.Execute(); err != nil {
		log.Fatal("failed to run command", err)
	}
}
