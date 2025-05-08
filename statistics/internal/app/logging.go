package app

import (
	"log"
	"os"
)

// InitializeLogging sets up logging for the application
func InitializeLogging() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Println("Logging initialized")
}