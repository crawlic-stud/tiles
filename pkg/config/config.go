package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PgDsn string
}

func Load() Config {
	// Ignore error in production.
	// Environment variables may already exist.
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	return Config{
		PgDsn: mustGetEnv("PG_DSN"),
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}
