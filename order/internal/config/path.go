package config

import (
	"os"
	"path/filepath"
)

// GetConfigPath returns the path to the configuration file
// It first checks for a CONFIG_PATH environment variable,
// and falls back to config/config.yaml in the working directory
func GetConfigPath() string {
	// Try to get config path from environment variable
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		return configPath
	}

	// Default to config/config.yaml in the working directory
	return filepath.Join("config", "config.yaml")
}
