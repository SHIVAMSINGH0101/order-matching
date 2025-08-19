package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for our application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "wifiname"),
			DBName:   getEnv("DB_NAME", "ordersdb"),
		},
	}

	return config, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
