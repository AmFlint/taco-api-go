package config

import (
	"github.com/AmFlint/taco-api-go/routes"
	"github.com/AmFlint/taco-api-go/routes/tasks"
	"github.com/AmFlint/taco-api-go/routes/lists"
)

// Function in charge of setting up Application Routes
func (a *App) initializeRoutes() {
	// ---- General Endpoints ---- //

	// Health Endpoint, check whether service is up or down
	a.Router.HandleFunc("/health", routes.HealthIndexHandler).Methods("GET")
	a.Router.HandleFunc("/health/", routes.HealthIndexHandler).Methods("GET")

	// ---- List Management Endpoints ---- //
	listRouter := a.Router.PathPrefix("/boards/{boardId}/lists").Subrouter()
	lists.InitRoutes(listRouter)

	// ---- Tasks Management Endpoints ---- //
	taskRouter:= listRouter.PathPrefix("/{listId}/tasks").Subrouter()
	tasks.InitRoutes(taskRouter)
}
