package config

import (
	"os"
)

// GetEnv returns the value of the environment variable or def if not set.
func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
