package main

import (
	"flag"
	"log"

	"github.com/soumikc1729/splitty/server/internal/config"
)

func main() {
	var configName, configPath string
	flag.StringVar(&configName, "config-name", "config", "config name")
	flag.StringVar(&configPath, "config-path", ".", "config path")

	cfg, err := config.Load[Config](configName, configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app, err := NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := app.Serve(); err != nil {
		app.Logger.Err(err).Msg("failed to start server")
	}
}
