package utils

import (
	"log"
	"os"
)

func GetAppURL() string {
	appURL := os.Getenv("APP_URL")

	if appURL == "" {
		log.Println("Warning: APP_URL is not set")
	}

	return appURL
}
