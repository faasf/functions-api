package main

import (
	"log"

	"github.com/faasf/functions-api/internal/app"
	"github.com/faasf/functions-api/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
