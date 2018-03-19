package main

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"net/http"
	"log"
	"io"
	"fmt"
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
}

// Retrieve Database connection from Application's Mongo Session
func (a *App) Connect() *mgo.Database {
	return a.DbSession.DB(databaseName)
}

// Parameter addr of form ":8080", represents the port where the application will be served
func (a *App) Run(addr string) {
	// Close Database connection at the end of Application Runtime
	defer a.DbSession.Close()

	// Routing policies takes place in this function
	a.setRoutes()

	log.Printf("Server listening on port%s", addr)
	// Listen on port defined in addr parameter, and serve Application via Mux Router
	// Configure Http server to log every access/error logs to Stdout
	if err := http.ListenAndServe(addr, handlers.LoggingHandler( os.Stdout, a.Router)); err != nil {
		log.Fatal(err)
	}
}

// Function in charge of setting up Application Routes
func (a *App) setRoutes() {
	// ---- General Endpoints ---- //

	// Health Endpoint, check whether service is up or down
	a.Router.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-type", "application/json")
		io.WriteString(writer, `{"alive": true}`)
	})

	a.Router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		hello := []byte("Hello, World")
		w.Write(hello)
	})

	// ---- Tasks Management Endpoints ---- //
	taskRouter := a.Router.PathPrefix("/boards/{boardId}/tasks").Subrouter()

	taskRouter.HandleFunc("/", routes.TaskIndexAction).Methods("GET")
}