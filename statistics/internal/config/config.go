package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	Kafka  KafkaConfig  `yaml:"kafka"`
	Logging LoggingConfig `yaml:"logging"`
}

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type DBConfig struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Name           string `yaml:"name"`
	SSLMode        string `yaml:"ssl_mode"`
	MaxConnections int    `yaml:"max_connections"`
}

type KafkaConfig struct {
	BootstrapServers string       `yaml:"bootstrap_servers"`
	ConsumerGroupID  string       `yaml:"consumer_group_id"`
	Topics           TopicsConfig `yaml:"topics"`
	AutoOffsetReset  string       `yaml:"auto_offset_reset"`
}

type TopicsConfig struct {
	OrderEvents   string `yaml:"order_events"`
	ProductEvents string `yaml:"product_events"`
	UserEvents    string `yaml:"user_events"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

func Load(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}