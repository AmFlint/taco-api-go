package config

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/gorilla/handlers"
	"os"
	"github.com/AmFlint/taco-api-go/config/database"
)

type App struct {
	Router    *mux.Router
}

var a App

// Create a new App structure from given configuration (Database Connection) + Initialize Routing
func NewApp(user, password, dbName, dbHost, dbPort string) App {
	a = App{}
	a.Initialize(user, password, dbName, dbHost, dbPort)
	return a
}

// Get Application Structure Instance from outside Pkgs
func GetApp() App {
	return a
}

// Initialize Application Structure: Configure Database Connection, save dbName for later uses and create Router
func (a *App) Initialize(user, password, dbName, dbHost, dbPort string) {

	database.SetDBSession(user, password, dbName, dbHost, dbPort)

	// Initialize Mux Router and assign it to application Structure
	a.Router = mux.NewRouter()
	// Routing policies takes place in this function
	// Routes are declared in -> routing.go
	a.initializeRoutes()
}

// Parameter addr of form ":8080", represents the port where the application will be served
func (a *App) Run(addr string) {
	// Close Database connection at the end of Application Runtime
	defer database.CloseSession()

	log.Printf("Server listening on port%s", addr)
	// Listen on port defined in addr parameter, and serve Application via Mux Router
	// Configure Http server to log every access/error logs to Stdout
	if err := http.ListenAndServe(addr, handlers.LoggingHandler( os.Stdout, a.Router)); err != nil {
		log.Fatal(err)
	}
}