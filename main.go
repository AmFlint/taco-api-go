package main

import (
	"os"
	"fmt"
)

// Custom Function to manage environment variables and use default values if doesn't exist
func getenv(key, fallback string) string {
	value := os.Getenv(key)
	// If Env Variable not defined, use fallback
	if len(value) == 0 {
		value = fallback
	}
	return value
}

// Application entry point, create Application Configuration with Environment Variables, generate Router + DB connection
// And serve the application
func main() {
	a := App{}
	a.Initialize(
		getenv("APP_USERNAME", ""),
		getenv("APP_PASSWORD", ""),
		getenv("APP_DB_NAME", "taco"),
		getenv("APP_DB_HOST", "localhost"),
		getenv("APP_DB_PORT", "27017"))

	// Get Port from Environment Variables and start server (Listen on given Port)
	port := fmt.Sprintf(":%s", getenv("APP_PORT", "8080"))
	a.Run(port)

}
