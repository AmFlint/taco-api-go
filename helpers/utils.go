package helpers

import (
	"net/http"
	"encoding/json"
	"fmt"
)

type ErrorMessages struct {
	Messages []string `json:"errors"`
	Code     int      `json:"code"`
}

const (
	ERROR__INVALID_PLAYLOAD = "Invalid Request Payload"
)

// ---- Write Http response from given code and payload ---- //
func writeResponse(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}

// ---- Write Http response with error message and status code ---- //
func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := fmt.Sprintf(`{"status": %v, "message": "%s"}`, code, message)
	writeResponse(w, code, []byte(response))
}

// ---- Write Http response With multiple error messages and status code ---- //
func RespondWithErrors(w http.ResponseWriter, code int, messages []string) {
	payload := ErrorMessages{Messages: messages, Code: code,}
	response, err := json.Marshal(payload)

	if err != nil {
		panic(err)
	}
	writeResponse(w, code, response)
}

// ---- Write Http response with json data, set header content type to json and set status code ---- //
func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	writeResponse(w, code, response)
}