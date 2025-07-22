
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
