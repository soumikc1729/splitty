package db

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Config struct {
	DSN          string        `mapstructure:"dsn"`
	QueryTimeout time.Duration `mapstructure:"query-timeout"`
	MaxOpenConns int           `mapstructure:"max-open-conns"`
	MaxIdleConns int           `mapstructure:"max-idle-conns"`
	IdleTimeout  time.Duration `mapstructure:"idle-timeout"`
	PingTimeout  time.Duration `mapstructure:"ping-timeout"`
}

type Models struct {
	Groups GroupModel
}

type Data struct {
	Config *Config
	Models Models
}

func New(cfg *Config) (*Data, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}

	data := Data{
		Config: cfg,
		Models: Models{
			Groups: GroupModel{DB: db},
		},
	}

	return &data, nil
}

func openDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.IdleTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
