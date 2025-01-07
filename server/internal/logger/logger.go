package logger

import (
	"errors"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type LogWriterType string

const (
	Console LogWriterType = "console"
)

var (
	ErrUnknownLogWriterType error = errors.New("unknown log writer type")
)

type Config struct {
	Level   string          `mapstructure:"level"`
	Writers []LogWriterType `mapstructure:"writers"`
}

func New(cfg *Config) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	zerolog.SetGlobalLevel(level)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var writers []io.Writer
	for _, writer := range cfg.Writers {
		switch writer {
		case Console:
			writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
		default:
			return nil, ErrUnknownLogWriterType
		}
	}

	logger := zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger()

	return &logger, nil
}
