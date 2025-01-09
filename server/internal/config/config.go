package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

func Load[T any](configName, configPath string) (*T, error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg T
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("error processing env variables: %v", err)
	}

	return &cfg, nil
}
