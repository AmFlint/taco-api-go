package testconfig

import (
	"github.com/AmFlint/taco-api-go/config"
	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/config/database"
	"os"
	"testing"
)

var app config.App

func clearDatabase() {
	// Retrieve Database Connection from main.Application Session
	db := database.GetDatabaseConnection()
	// Clear Database Content
	db.DropDatabase()
}

func Init(m *testing.M) {
	// Create main.Application Struct
	// Initialize main.Application with database configuration
	if app == (config.App{}) {
		app = config.NewApp(
			helpers.GetEnv("APP_USERNAME", ""),
			helpers.GetEnv("APP_PASSWORD", ""),
			helpers.GetEnv("APP_DB_NAME", "taco"),
			helpers.GetEnv("APP_DB_HOST", "localhost"),
			helpers.GetEnv("APP_DB_PORT", "27017"))
	}
	
	// Clear database collections
	clearDatabase()

	// Run following tests
	code := m.Run()

	// Clear Database after use
	clearDatabase()

	defer database.CloseSession()

	// Exit testing program with runtime exit code
	os.Exit(code)
}