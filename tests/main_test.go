package tests

import (
	"testing"
	"os"
	"github.com/AmFlint/taco-api-go/config"
	"github.com/AmFlint/taco-api-go/helpers"
)

var a config.App

// Clear Database Collections
func clearDatabase(a config.App) {
	// Retrieve Database Connection from main.Application Session
	db := a.Connect()
	// Clear Database Content
	db.DropDatabase()
}

func TestMain(m *testing.M) {
	// Create main.Application Struct
	// Initialize main.Application with database configuration
	a := config.NewApp(
		helpers.GetEnv("main.App_USERNAME", ""),
		helpers.GetEnv("main.App_PASSWORD", ""),
		helpers.GetEnv("main.App_DB_NAME", "taco"),
		helpers.GetEnv("main.App_DB_HOST", "localhost"),
		helpers.GetEnv("main.App_DB_PORT", "27017"))

	// Clear database collections
	clearDatabase(a)

	// Run following tests
	code := m.Run()

	// Clear Database after use
	clearDatabase(a)

	defer a.DbSession.Close()

	// Exit testing program with runtime exit code
	os.Exit(code)
}
