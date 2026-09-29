package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	ServerPort       string
	DatabaseURL      string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	FrontendURL      string
}

func LoadConfig() *Config {
	port := getEnv("SERVER_PORT", "8080")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "supermarket")
	dbSSL := getEnv("DB_SSLMODE", "disable")

	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, dbSSL)
	}

	accessExpStr := getEnv("JWT_ACCESS_EXPIRY", "15m")
	accessExp, err := time.ParseDuration(accessExpStr)
	if err != nil {
		accessExp = 15 * time.Minute
	}

	refreshExpStr := getEnv("JWT_REFRESH_EXPIRY", "168h")
	refreshExp, err := time.ParseDuration(refreshExpStr)
	if err != nil {
		refreshExp = 168 * time.Hour
	}

	return &Config{
		ServerPort:       port,
		DatabaseURL:      dbURL,
		JWTSecret:        getEnv("JWT_SECRET", "super-secret-key-change-me"),
		JWTAccessExpiry:  accessExp,
		JWTRefreshExpiry: refreshExp,
		FrontendURL:      getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
