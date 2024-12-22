package main

import (
	"log"

	"forum/internal/app"
	"forum/pkg/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
