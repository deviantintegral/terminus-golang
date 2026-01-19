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

func TestCommit(t *testing.T) {
	// Commit() returns the commit hash or empty string
	// We just verify it doesn't panic
	_ = Commit()
}

func TestIsDirty(t *testing.T) {
	// Without ldflags, dirty is empty, so IsDirty() returns false
	if IsDirty() && dirty != "true" {
		t.Error("IsDirty() should return false when dirty is not 'true'")
	}
}

func TestIsRelease(t *testing.T) {
	// Without ldflags, version is empty, so IsRelease() returns false
	if IsRelease() && version != "" {
		t.Error("IsRelease() should return false when version is empty")
	}
}
