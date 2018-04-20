package tasks

import (
	"github.com/gorilla/mux"
)

// Initialize Routes for Task Resource
func InitRoutes(taskRouter *mux.Router) {
	taskRouter.HandleFunc("", TaskIndexHandler).Methods("GET")
	taskRouter.HandleFunc("", TaskCreateHandler).Methods("POST")
	taskRouter.HandleFunc("/{taskId}", TaskViewHandler).Methods("GET")
	taskRouter.HandleFunc("/{taskId}", TaskDeleteHandler).Methods("DELETE")
}