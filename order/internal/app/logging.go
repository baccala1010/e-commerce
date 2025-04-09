package app

import (
	"time"

	"github.com/baccala1010/e-commerce/order/internal/config"
	"github.com/sirupsen/logrus"
)

// SetupLogging configures the global logger based on application configuration
func SetupLogging(cfg *config.Config) {
	// Set log level based on configuration
	level := logrus.InfoLevel
	if cfg.Logging.Level == "debug" {
		level = logrus.DebugLevel
	}
	logrus.SetLevel(level)

	// Set log formatter
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

	logrus.Infof("Log level set to %s", cfg.Logging.Level)
}
