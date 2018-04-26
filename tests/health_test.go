package tests

import (
	"testing"
	"net/http"
	"github.com/AmFlint/taco-api-go/tests/utils"
	"github.com/AmFlint/taco-api-go/tests/utils/testconfig"
)

func TestMain(m *testing.M) {
	testconfig.Init(m)
}

// Test that Health Check endpoint on "/health" returns a healthy response
func TestHealthCheck(t *testing.T) {
	// Create GET Http Request to endpoint /health with no body
	req, _ := http.NewRequest("GET", "/health", nil)
	// Execute Request and retrieve response
	response := utils.ExecuteRequest(req)

	// Check that Response Code is 200 | OK
	utils.CheckResponseCode(t, response.Code, http.StatusOK)
	// Assert that response body is {"alive": true}
	utils.AssertStringEqualsTo(t, response.Body.String(), `{"alive": true}`)
}