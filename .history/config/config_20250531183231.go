package config

import (
	"os"
	"strconv"
)


type Config struct {
	AppPort  int
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
}


func GetConfig() *Config {
	config := &Config{
		AppPort:  8081,
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASS", "postgres"),
		DBName:   getEnv("DB_NAME", "premier_league"),
	}
	

	portStr := getEnv("APP_PORT", "8081")
	port, err := strconv.Atoi(portStr)
	if err == nil {
		config.AppPort = port
	}
	
	return config
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 