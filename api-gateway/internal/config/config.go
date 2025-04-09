package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port int
	Name string
}

type ServicesConfig struct {
	Inventory ServiceConfig
	Order     ServiceConfig
}

type ServiceConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

type LoggingConfig struct {
	Level string
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
