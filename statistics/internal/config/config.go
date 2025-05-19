package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPC struct {
		Port int `yaml:"port"`
	} `yaml:"grpc"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		GroupID string   `yaml:"group_id"`
		Topics  []string `yaml:"topics"`
	} `yaml:"kafka"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
