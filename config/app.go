package config

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/gorilla/handlers"
	"os"
	"github.com/AmFlint/taco-api-go/routes"
)

type App struct {
	Router    *mux.Router
	DbSession *mgo.Session
}

// Variable to contain the database name for later use when Session connects to Database
var databaseName string
var a App

func NewApp(user, password, dbName, dbHost, dbPort string) App {
	a = App{}
	a.Initialize(user, password, dbName, dbHost, dbPort)
	return a
}

func GetApp() App {
	return a
}

// Initialize Application Structure: Configure Database Connection, save dbName for later uses and create Router
func (a *App) Initialize(user, password, dbName, dbHost, dbPort string) {

	databaseName = dbName
	var dbUrl string
	// Configure dbUrl, manage user/password credentials, if provided
	if len(user) > 0 && len(password) > 0 {
		dbUrl += fmt.Sprintf("%s@%s", user, password)
	}
	// Add Mongo Host:Port to connection url
	dbUrl += fmt.Sprintf("%s:%s", dbHost, dbPort)

	var err error
	// Connect to MongoDb server and retrieve Session Instance
	a.DbSession, err = mgo.Dial(dbUrl)
	if err != nil {
		panic(err)
	}

	log.Print("Connection to Database established with success !")

	// Initialize Mux Router and assign it to application Structure
	a.Router = mux.NewRouter()
	// Routing policies takes place in this function
	a.initializeRoutes()
}

// Retrieve Database connection from Application's Mongo Session
func (a *App) Connect() *mgo.Database {
	return a.DbSession.DB(databaseName)
}

// Parameter addr of form ":8080", represents the port where the application will be served
func (a *App) Run(addr string) {
	// Close Database connection at the end of Application Runtime
	defer a.DbSession.Close()

	log.Printf("Server listening on port%s", addr)
	// Listen on port defined in addr parameter, and serve Application via Mux Router
	// Configure Http server to log every access/error logs to Stdout
	if err := http.ListenAndServe(addr, handlers.LoggingHandler( os.Stdout, a.Router)); err != nil {
		log.Fatal(err)
	}
}

// Function in charge of setting up Application Routes
func (a *App) initializeRoutes() {
	// ---- General Endpoints ---- //

	// Health Endpoint, check whether service is up or down
	a.Router.HandleFunc("/health", routes.HealthIndexHandler).Methods("GET")

	// ---- Tasks Management Endpoints ---- //
	//taskRouter := a.Router.PathPrefix("/boards/{boardId}/tasks").Subrouter()

	//taskRouter.HandleFunc("/", routes.TaskIndexAction).Methods("GET")
}