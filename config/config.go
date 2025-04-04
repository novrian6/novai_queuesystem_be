package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct for storing configuration values
type Config struct {
	DatabaseURL string
	JWTSecret   string
	SMTPServer  string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
}

// LoadConfig reads configuration from .env file and environment variables
func LoadConfig() *Config {
	// Load .env file if available
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/queue_system"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secure-secret-key"),
		SMTPServer:  getEnv("SMTP_SERVER", "smtp.example.com"),
		SMTPPort:    getEnv("SMTP_PORT", "587"),
		SMTPUser:    getEnv("SMTP_USER", "no-reply@example.com"),
		SMTPPass:    getEnv("SMTP_PASSWORD", ""),
	}
}

// getEnv retrieves environment variables or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
