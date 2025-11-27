package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	AllowedOrigins []string
	Environment    string
}

// LoadConfig reads .env and returns config
// SIMPLE & STRICT: All required values must be set
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get all required values
	databaseURL := getEnvRequired("DATABASE_URL")
	jwtSecret := getEnvRequired("JWT_SECRET")
	allowedOriginsStr := getEnvRequired("ALLOWED_ORIGINS")
	
	// Optional with defaults
	port := getEnvDefault("PORT", "8080")
	environment := getEnvDefault("ENV", "development")

	// Parse ALLOWED_ORIGINS
	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	log.Printf("✓ Config loaded - Environment: %s, Port: %s", environment, port)

	return &Config{
		DatabaseURL:    databaseURL,
		JWTSecret:      jwtSecret,
		Port:           port,
		AllowedOrigins: allowedOrigins,
		Environment:    environment,
	}
}

// getEnvRequired gets env var or fails
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("ERROR: %s is required in .env", key)
	}
	return value
}

// getEnvDefault gets env var or returns default
func getEnvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
