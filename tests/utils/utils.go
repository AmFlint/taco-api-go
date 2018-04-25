package utils

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/AmFlint/taco-api-go/config"
	"reflect"
	"fmt"
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
// TODO: Use func assertEquals in all AssertEquals methods
func assertEquals(t *testing.T, got, expected interface{}, err string) {
	if got != expected {
		t.Error(err)
	}
}

func AssertStringEqualsTo(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("[Error], Expected body to be: %s, got: %s", expected, got)
	}
}

func AssertBoolEqualsTo(t *testing.T, got, expected bool) {
	assertEquals(t, got, expected, fmt.Sprintf("[Error], Expected given boolean to be %v, got %v", expected, got))
}

func AssertFloatEqualsTo(t *testing.T, got, expected float64) {
	if got != expected {
		t.Errorf("[Error], Expected value to be: %v, got: %v", expected, got)
	}
}

func AssertIntEqualsTo(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("[Error], Expected integer value to be: %v, got: %v", expected, got)
	}
}

func AssertMapHasKey(t *testing.T, m map[string]interface{}, key string) {
	if m[key] == nil {
		t.Errorf("[Error], Expected map to have key %s, but does not exist", key)
	}
}

// Assert that interface is not Empty
func AssertNotEmpty(t *testing.T, got interface{}) {
	if reflect.DeepEqual(got, reflect.Zero(reflect.TypeOf(got)).Interface()) {
		t.Errorf("[Error] Expected interface to be non-empty, but is empty")
	}
}