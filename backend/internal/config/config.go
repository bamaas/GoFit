package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type Config struct {
	LogLevel string
	DevelopmentMode bool
	Users []user
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

type user struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func Retrieve() (*Config, error){

	// Get development mode
	devMode, err := strconv.ParseBool((getEnv("DEVELOPMENT_MODE", defaultDevelopmentMode)))
	if err != nil {
		return nil, err
	}
 
	// Get users
	var users []user
	err = json.Unmarshal([]byte(getEnv("USERS", `[]`)), &users)
	if err != nil {
		return nil, err
	}

	return &Config{
		LogLevel: getEnv("LOG_LEVEL", defaultLogLevel),
		DevelopmentMode: devMode,
		Users: users,
	}, nil
}
