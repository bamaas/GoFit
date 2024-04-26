package config

import (
	"os"
)

type Config struct {
	LogLevel string
}

// Default configuration values
var defaultLogLevel = "INFO"

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func Retrieve() (*Config, error){
	return &Config{
		LogLevel: getEnv("LOG_LEVEL", defaultLogLevel),
	}, nil
}
