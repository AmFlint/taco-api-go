package main

import (
	"fmt"
	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/config"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Setting-up logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

// Application entry point, create Application Configuration with Environment Variables, generate Router + DB connection
// And serve the application
func main() {
	a := config.NewApp(
		helpers.GetEnv("APP_USERNAME", ""),
		helpers.GetEnv("APP_PASSWORD", ""),
		helpers.GetEnv("APP_DB_NAME", "taco"),
		helpers.GetEnv("APP_DB_HOST", "localhost"),
		helpers.GetEnv("APP_DB_PORT", "27017"))

	// Get Port from Environment Variables and start server (Listen on given Port)
	port := fmt.Sprintf(":%s", helpers.GetEnv("APP_PORT", "8080"))
	a.Run(port)
}
