package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func Load() *Config {

	if os.Getenv("APP_ENV") == "DEV" {
		if err := godotenv.Load(".env.dev"); err != nil {
			log.Fatal("Error loading .env.dev file")
		}
	}
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
