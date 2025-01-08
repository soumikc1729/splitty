package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Load[T any](configName, configPath string) (*T, error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg T
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &cfg, nil
}
