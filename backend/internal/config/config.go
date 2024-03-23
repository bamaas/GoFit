package config

import (
	"os"
)

type config struct {
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

func Get() (*config, error){

	return &config{
		LogLevel: getEnv("LOG_LEVEL", defaultLogLevel),
	}, nil
}
