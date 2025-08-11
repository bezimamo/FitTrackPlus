package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
// In Go, we use structs to group related data
type Config struct {
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	
	// Server configuration
	Port string
	
	// JWT configuration
	JWTSecret string
}

// LoadConfig loads configuration from environment variables
// This function returns a pointer to Config struct
func LoadConfig() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	return &Config{
		// Database settings
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "fittrackplus"),
		
		// Server settings
		Port: getEnv("PORT", "8080"),
		
		// JWT settings
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}

// getEnv is a helper function to get environment variables with defaults
// In Go, we can return multiple values from a function
func getEnv(key, defaultValue string) string {
	// os.Getenv gets the value of an environment variable
	value := os.Getenv(key)
	
	// If the environment variable is empty, return the default value
	if value == "" {
		return defaultValue
	}
	
	return value
} 