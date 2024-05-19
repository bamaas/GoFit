package config

import (
	"os"
	"strconv"
)

type Config struct {
	LogLevel string
	DevelopmentMode bool
}

// Default configuration values
var defaultLogLevel = "INFO"
var defaultDevelopmentMode = "false"

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}	

func Retrieve() (*Config, error){

	// Get development mode
	devMode, err := strconv.ParseBool((getEnv("DEVELOPMENT_MODE", defaultDevelopmentMode)))
	if err != nil {
		return nil, err
	}

	return &Config{
		LogLevel: getEnv("LOG_LEVEL", defaultLogLevel),
		DevelopmentMode: devMode,
	}, nil
}
