package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL    string
	ExternalAPIURL string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	config := &Config{
		PostgresURL:    os.Getenv("POSTGRES_URL"),
		ExternalAPIURL: os.Getenv("EXTERNAL_API_URL"),
	}

	return config, nil
}
