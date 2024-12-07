package config

import (
	"os"

	"github.com/joho/godotenv"
)

func GetDatabaseURL() string {
	godotenv.Load()

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		DATABASE_URL = "postgres://admin:admin@localhost:5432/db?sslmode=disable"
	}
	return DATABASE_URL
}

func GetPort() string {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	return PORT
}

func GetSecretKey() string {
	godotenv.Load()

	SECRET_KEY := os.Getenv("SECRET_KEY")
	if SECRET_KEY == "" {
		SECRET_KEY = "secret"
	}
	return SECRET_KEY
}

func GetRedisEndpoint() string {
	godotenv.Load()

	REDIS_ENDPOINT := os.Getenv("REDIS_URL")
	if REDIS_ENDPOINT == "" {
		REDIS_ENDPOINT = "localhost:6379"
	}
	return REDIS_ENDPOINT
}

func GetRedisPassword() string {
	godotenv.Load()

	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	if REDIS_PASSWORD == "" {
		REDIS_PASSWORD = ""
	}
	return REDIS_PASSWORD
}
