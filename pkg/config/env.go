package config

import (
	"enlabs-task/pkg/model"
	"os"
	"strconv"
	"time"
)

// NewConfig creates a new Config instance with values from environment variables
func NewConfig() *model.Config {
	return &model.Config{
		Server: model.Server{
			Port:         getEnv("API_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},

		Database: model.Database{
			Host:               getEnv("DB_HOST", "localhost"),
			Port:               getEnv("DB_PORT", "5432"),
			User:               getEnv("DB_USER", "postgres"),
			Password:           getEnv("DB_PASSWORD", "postgres"),
			DBName:             getEnv("DB_NAME", "transactions"),
			SSLMode:            getEnv("DB_SSLMODE", "disable"),
			MaxOpenConnections: getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConnections: getIntEnv("DB_MAX_IDLE_CONNS", 5),
		},

		App: model.App{
			Port:        getEnv("APP_PORT", "8080"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
	}
}

// Helper functions for getting environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
