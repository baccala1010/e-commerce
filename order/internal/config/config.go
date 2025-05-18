package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server           ServerConfig
	Database         DatabaseConfig
	InventoryService InventoryServiceConfig `mapstructure:"inventory_service"`
	Kafka            KafkaConfig
	Logging          LoggingConfig
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

type InventoryServiceConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

type LoggingConfig struct {
	Level string
}

type KafkaConfig struct {
	BootstrapServers string `mapstructure:"bootstrap_servers"`
	Topics           KafkaTopics
}

type KafkaTopics struct {
	OrderEvents string `mapstructure:"order_events"`
	UserEvents  string `mapstructure:"user_events"`
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
