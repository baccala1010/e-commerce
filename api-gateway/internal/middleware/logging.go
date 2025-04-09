package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger middleware for logging HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		requestPath := c.Request.URL.Path

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request status
		status := c.Writer.Status()

		// Set log level based on status code
		var logFunc func(args ...interface{})
		if status >= 500 {
			logFunc = logrus.Error
		} else if status >= 400 {
			logFunc = logrus.Warn
		} else {
			logFunc = logrus.Info
		}

		// Log request details
		logFunc(map[string]interface{}{
			"status":     status,
			"method":     c.Request.Method,
			"path":       requestPath,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
		})
	}
}
