package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort      string
	LogLevel        string
	DefaultPageSize int
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
}

func LoadConfig() Config {
	return Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		DefaultPageSize: getEnvAsInt("DEFAULT_PAGE_SIZE", 10),
		DBHost:          getEnv("DB_HOST", "postgres"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "taskdb"),
		DBSSLMode:       getEnv("DB_SSL_MODE", "disable"),
	}
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists && valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil && value > 0 {
			return value
		}
	}
	return defaultValue
}
