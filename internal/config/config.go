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

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("[DEBUG] Environment URL: ", os.Getenv("DATABASE_URL"))
	log.Println("[DEBUG] Database URL: ", os.Getenv("DATABASE_URL"))
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
