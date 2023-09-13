package main

import (
	"log"

	"github.com/labasubagia/realworld-backend/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal("failed to run command", err)
	}
}
