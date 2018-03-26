package database

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"log"
)

// Variables used to hold Servers's Database Session and database name for later database connections
var DbSession *mgo.Session
var databaseName string

// Return DatabaseSession
func GetDBSession() *mgo.Session {
	return DbSession
}

// Initialize Database Session with Mongo server according to given parameters (credentials, mongo address)
func SetDBSession(user, password, dbName, dbHost, dbPort string) {
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
	DbSession, err = mgo.Dial(dbUrl)
	if err != nil {
		panic(err)
	}

	log.Print("Connection to Database established with success !")
}

// Return Database Connection
func GetDatabaseConnection() *mgo.Database {
	return DbSession.DB(databaseName)
}

// Close Db Session - Unlink App -> Mongo server
func CloseSession() {
	DbSession.Close()
}