package routes

import (
	"net/http"
	"io"
)

// Health endpoint on GET /health, returns http status 200 ans payload {"alive": true}, useful for application health-checks
func HealthIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}