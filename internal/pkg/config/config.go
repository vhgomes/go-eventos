package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	MaxOpenConns int
	MaxIdleConns int
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required (e.g., postgres://user:pass@localhost:5432/db?sslmode=disable)")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required (min 32 chars)")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}

	maxOpen := getEnvInt("DB_MAX_OPEN_CONNS", 25)
	maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", 5)

	return &Config{
		Port:         port,
		DatabaseURL:  databaseURL,
		JWTSecret:    jwtSecret,
		MaxOpenConns: maxOpen,
		MaxIdleConns: maxIdle,
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
