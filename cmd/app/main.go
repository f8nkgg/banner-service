package main

import (
	"banner/config"
	"banner/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
