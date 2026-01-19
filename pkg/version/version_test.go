package version

import (
	"testing"
)

func TestString(t *testing.T) {
	// When no build-time values are set, String() returns "dev"
	// This is the expected behavior for local development builds
	result := String()
	if result == "" {
		t.Error("String() should never return an empty string")
	}
	// Default value without ldflags should be "dev"
	if result != "dev" {
		t.Logf("String() returned %q (may have build-time values set)", result)
	}
}
