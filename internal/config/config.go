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

type Config struct {
	DatabaseURL string
}

func Load() {
	if os.Getenv("ENVIRONMENT") == "PROD" {
		loadFromSecretsManager()
	} else {
		loadFromDotEnv()
	}
}

func loadFromDotEnv() {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Fatal("Error loading .env.dev file: ", err)
	}
}

func loadFromSecretsManager() {
	secretName := "ProdEnv"
	region := "ap-southeast-1"

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
		if key != "SSL_CERT_DATA" && key != "SSL_KEY_DATA" {
			log.Printf("[INFO] Environment variable %s: %s", key, os.Getenv(key))
		}
	}
}
