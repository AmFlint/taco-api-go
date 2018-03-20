package helpers

import "os"

// Custom Function to manage environment variables and use default values if doesn't exist
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	// If Env Variable not defined, use fallback
	if len(value) == 0 {
		value = fallback
	}
	return value
}
