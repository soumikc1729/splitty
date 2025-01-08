package main

import (
	"github.com/rs/zerolog"
	"github.com/soumikc1729/splitty/server/internal/db"
	"github.com/soumikc1729/splitty/server/internal/logger"
	"github.com/soumikc1729/splitty/server/internal/server"
)

type Config struct {
	Server server.Config `mapstructure:"server"`
	Logger logger.Config `mapstructure:"logger"`
	DB     db.Config     `mapstructure:"db"`
}

type App struct {
	Config *Config
	Logger *zerolog.Logger
}

func NewApp(cfg *Config) (*App, error) {
	logger, err := logger.New(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &App{Config: cfg, Logger: logger}, nil
}

func (app *App) Serve() error {
	server := server.New(&app.Config.Server, app.Logger)
	return server.Start(app.Routes())
}
