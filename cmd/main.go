package main

import (
	"L0/config"
	"L0/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("error at start: %v", err)
	}

	app.Run(cfg)
}
