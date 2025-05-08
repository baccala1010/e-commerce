package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetPath returns the path to the config file based on the environment
func GetPath(fileName string) (string, error) {
	isDocker := os.Getenv("DOCKER") == "true"

	var configFileName string
	if isDocker {
		configFileName = "config.docker.yaml"
	} else {
		configFileName = "config.yaml"
	}

	if fileName != "" {
		configFileName = fileName
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Try to find config directory up to 3 levels up
	for i := 0; i <= 3; i++ {
		configPath := filepath.Join(wd, "config", configFileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
		wd = filepath.Dir(wd)
	}

	return "", fmt.Errorf("config file %s not found", configFileName)
}