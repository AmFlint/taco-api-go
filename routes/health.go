package routes

import (
	"net/http"
	"io"
)

func HealthIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}