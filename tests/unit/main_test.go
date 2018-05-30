package unit

import (
	"testing"
	"github.com/AmFlint/taco-api-go/tests/utils/testconfig"
)

// TestMain - Initialize application for Testing
func TestMain(m *testing.M) {
	testconfig.Init(m)
}