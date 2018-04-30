package lists

import (
	"github.com/gorilla/mux"
)

// Initialize Routes for Task Resource
func InitRoutes(listRouter *mux.Router) {
	// ---- List Creation ---- //
	listRouter.HandleFunc("", ListCreateHandler).Methods("POST")
	listRouter.HandleFunc("/", ListCreateHandler).Methods("POST")
	// ---- List Deletion ---- //
	listRouter.HandleFunc("/{listId}", ListDeleteHandler).Methods("DELETE")
	listRouter.HandleFunc("/{listId}/", ListDeleteHandler).Methods("DELETE")
	// ---- List View ---- //
	listRouter.HandleFunc("/{listId}", ListViewHandler).Methods("GET")
	listRouter.HandleFunc("/{listId}/", ListViewHandler).Methods("GET")
}