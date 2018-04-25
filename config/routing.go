package config

import (
	"github.com/AmFlint/taco-api-go/routes"
	"github.com/AmFlint/taco-api-go/routes/tasks"
)

// Function in charge of setting up Application Routes
func (a *App) initializeRoutes() {
	// ---- General Endpoints ---- //

	// Health Endpoint, check whether service is up or down
	a.Router.HandleFunc("/health", routes.HealthIndexHandler).Methods("GET")
	a.Router.HandleFunc("/health/", routes.HealthIndexHandler).Methods("GET")

	// ---- Tasks Management Endpoints ---- //
	taskRouter := a.Router.PathPrefix("/boards/{boardId}/lists/{listId}/tasks").Subrouter()

	tasks.InitRoutes(taskRouter)
}
