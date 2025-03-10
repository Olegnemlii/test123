package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config stores the application configuration
type Config struct {
	Port            string
	DatabaseURL     string
	MailopostApiKey string
	MailopostURL    string
}

// LoadConfig loads the configuration from environment variables or .env file
func LoadConfig(path string) (*Config, error) {
	err := godotenv.Load(path + "/.env") // load .env file
	if err != nil {
		log.Printf("could not load .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051" // default port
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	mailopostApiKey := os.Getenv("MAILOPOST_API_KEY")
	if mailopostApiKey == "" {
		return nil, fmt.Errorf("MAILOPOST_API_KEY is not set")
	}

	mailopostURL := os.Getenv("MAILOPOST_URL")
	if mailopostURL == "" {
		return nil, fmt.Errorf("MAILOPOST_URL is not set")
	}

	return &Config{
		Port:            port,
		DatabaseURL:     databaseURL,
		MailopostApiKey: mailopostApiKey,
		MailopostURL:    mailopostURL,
	}, nil
}
