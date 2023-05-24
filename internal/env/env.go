package env

import (
	"os"
	"strconv"
)

// GetStringEnvOrDefault returns the value of the environment variable key.
// If the environment variable is not set, it returns def.
func GetStringEnvOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

// GetIntEnvOrDefault returns the value of the environment variable key as an int.
// If the environment variable is not set or cannot be parsed as an int, it returns def.
func GetIntEnvOrDefault(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		ival, err := strconv.Atoi(val)
		if err != nil {
			return def
		}
		return ival
	}
	return def
}
