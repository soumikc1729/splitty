package main

import (
	"github.com/rs/zerolog"
	"github.com/soumikc1729/splitty/server/internal/data"
	"github.com/soumikc1729/splitty/server/internal/logger"
	"github.com/soumikc1729/splitty/server/internal/server"
)

type Config struct {
	Server server.Config `mapstructure:"server"`
	Logger logger.Config `mapstructure:"logger"`
	Data   data.Config   `mapstructure:"data"`
}

type App struct {
	Config *Config
	Logger *zerolog.Logger
	Data   *data.Data
}

func NewApp(cfg *Config) (*App, error) {
	logger, err := logger.New(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	data, err := data.New(&cfg.Data)
	if err != nil {
		return nil, err
	}

	return &App{Config: cfg, Logger: logger, Data: data}, nil
}

func (app *App) Serve() error {
	server := server.New(&app.Config.Server, app.Logger)
	return server.Start(app.Routes())
}
