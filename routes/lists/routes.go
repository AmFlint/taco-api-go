package lists

import (
	"github.com/gorilla/mux"
)

// Initialize Routes for Task Resource
func InitRoutes(listRouter *mux.Router) {
	// ---- List Creation ---- //
	listRouter.HandleFunc("", ListCreateHandler).Methods("POST")
	listRouter.HandleFunc("/", ListCreateHandler).Methods("POST")
}