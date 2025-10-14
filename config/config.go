package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://productuser:productpass@localhost:5432/productdb?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}