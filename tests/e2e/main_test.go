package e2e

import (
	"testing"
	"github.com/AmFlint/taco-api-go/tests/utils/testconfig"
)

// Initializing Application with Database Session (re-use same DB Connection in every testing suite)
func TestMain(m *testing.M) {
	testconfig.Init(m)
}
