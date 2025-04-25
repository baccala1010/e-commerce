package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port     int
	GRPCPort int `mapstructure:"grpc_port"`
	Name     string
}

type DatabaseConfig struct {
	Host                  string
	Port                  int
	Name                  string
	Username              string
	Password              string
	SSLMode               string `mapstructure:"sslmode"`
	MaxIdleConnections    int    `mapstructure:"max_idle_connections"`
	MaxOpenConnections    int    `mapstructure:"max_open_connections"`
	ConnectionMaxLifetime string `mapstructure:"connection_max_lifetime"`
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

func (dc *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dc.Host, dc.Port, dc.Username, dc.Password, dc.Name, dc.SSLMode)
}

func (dc *DatabaseConfig) GetConnectionMaxLifetime() time.Duration {
	d, err := time.ParseDuration(dc.ConnectionMaxLifetime)
	if err != nil {
		return time.Hour
	}
	return d
}
