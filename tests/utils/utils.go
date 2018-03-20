package utils

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/AmFlint/taco-api-go/config"
)

// Execute Http Request
func ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	a := config.GetApp()
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

// Asserts an error if "Got" Reponse Code is different than "Expected" Response code
func CheckResponseCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("Got Response Code %v, expected: %v", got, expected)
	}
}

func AssertResponseStringIs(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("[Error], Expected Response body to be: %s, got: %s", expected, got)
	}
}