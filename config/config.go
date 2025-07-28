package config

import (
	"os"
	"time"
)

type Config struct {
	Port string
	DBFile string
	JWTSecret string
	JWTExpiry time.Duration
	RefreshExpiry time.Duration
	AllowedOrigins []string
}


func LoadConfig() Config {
    return Config{
        Port:           getEnv("PORT", "8080"),
        DBFile:         getEnv("DB_FILE", "./secure-api.db"),
        JWTSecret:      getEnv("JWT_SECRET", "a-very-secret-key-that-is-long-and-random"),
        JWTExpiry:      15 * time.Minute,
        RefreshExpiry:  7 * 24 * time.Hour,
        AllowedOrigins: []string{"http://localhost:3000"},
    }
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback

}