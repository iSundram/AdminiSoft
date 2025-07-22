
package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DatabaseURL     string
	RedisURL        string
	JWTSecret       string
	Environment     string
	AdminEmail      string
	AdminPassword   string
	EncryptionKey   string
}

func Load() *Config {
	godotenv.Load()
	
	return &Config{
		Port:            getEnv("PORT", "5000"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://user:password@localhost/adminisoftware?sslmode=disable"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		AdminEmail:      getEnv("ADMIN_EMAIL", "admin@adminisoftware.com"),
		AdminPassword:   getEnv("ADMIN_PASSWORD", "admin123"),
		EncryptionKey:   getEnv("ENCRYPTION_KEY", "your-encryption-key-32-chars"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	Environment  string
	DatabaseURL  string
	RedisURL     string
	JWTSecret    string
	Debug        bool
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "5000"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://adminisoftware:adminisoftware123@localhost:5432/adminisoftware_db?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-this"),
		Debug:        getEnvAsBool("DEBUG", true),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
