package tasks

import (
	"github.com/gorilla/mux"
)

// Initialize Routes for Task Resource
func InitRoutes(taskRouter *mux.Router) {
	// ---- Task Listing ---- //
	taskRouter.HandleFunc("", TaskIndexHandler).Methods("GET")
	taskRouter.HandleFunc("/", TaskIndexHandler).Methods("GET")
	// ----  Task Creation ---- //
	taskRouter.HandleFunc("", TaskCreateHandler).Methods("POST")
	taskRouter.HandleFunc("/", TaskCreateHandler).Methods("POST")
	// ---- Task View  ---- //
	taskRouter.HandleFunc("/{taskId}", TaskViewHandler).Methods("GET")
	taskRouter.HandleFunc("/{taskId}/", TaskViewHandler).Methods("GET")
	// ---- Task Deletion ---- //
	taskRouter.HandleFunc("/{taskId}", TaskDeleteHandler).Methods("DELETE")
	taskRouter.HandleFunc("/{taskId}/", TaskDeleteHandler).Methods("DELETE")
}