package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	DatabaseURL string
}

func Load() *Config {
	if os.Getenv("ENVIRONMENT") == "PROD" {
		return loadFromSecretsManager()
	} else {
		return loadFromDotEnv()
	}
}

// loadFromDotEnv loads environment variables from a .env file
func loadFromDotEnv() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	log.Println("[DEBUG] Database URL: ", os.Getenv("DATABASE_URL"))
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

// loadFromSecretsManager fetches configuration from AWS Secrets Manager
func loadFromSecretsManager() *Config {
	secretName := "ProdEnv"    // Specify your secret name here
	region := "ap-southeast-1" // Specify the AWS region here

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal("Unable to load AWS SDK config: ", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal("Failed to retrieve secret from AWS Secrets Manager: ", err)
	}

	var secretData map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secretData); err != nil {
		log.Fatal("Failed to unmarshal secret string: ", err)
	}

	for key, value := range secretData {
		if err := os.Setenv(key, value); err != nil {
			log.Fatalf("Failed to set environment variable %s, %v", key, err)
		}
		log.Printf("[DEBUG] Set environment variable %s: %s", key, os.Getenv(key))
	}

	dbURL, exists := secretData["DATABASE_URL"]
	if !exists {
		log.Fatal("DATABASE_URL not found in secrets manager")
	}

	log.Println("[DEBUG] Database URL from AWS Secrets Manager: ", dbURL)
	return &Config{
		DatabaseURL: dbURL,
	}
}
