package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMultidevService_List(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sites/site1/environments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Return a mix of standard and multidev environments
		envs := []map[string]interface{}{
			{"id": "dev", "domain": "dev.site1.pantheonsite.io"},
			{"id": "test", "domain": "test.site1.pantheonsite.io"},
			{"id": "live", "domain": "live.site1.pantheonsite.io"},
			{"id": "feature-1", "domain": "feature-1.site1.pantheonsite.io"},
			{"id": "feature-2", "domain": "feature-2.site1.pantheonsite.io"},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(envs)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	multidevService := NewMultidevService(client)

	multidevs, err := multidevService.List(context.Background(), "site1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should only return multidev environments (not dev, test, live)
	if len(multidevs) != 2 {
		t.Errorf("expected 2 multidev environments, got %d", len(multidevs))
	}

	// Check that we got the correct environments
	foundFeature1 := false
	foundFeature2 := false
	for _, env := range multidevs {
		if env.ID == "feature-1" {
			foundFeature1 = true
		}
		if env.ID == "feature-2" {
			foundFeature2 = true
		}
		// Should not include standard environments
		if env.ID == "dev" || env.ID == "test" || env.ID == "live" {
			t.Errorf("should not include standard environment '%s'", env.ID)
		}
	}

	if !foundFeature1 {
		t.Error("expected to find 'feature-1' multidev")
	}
	if !foundFeature2 {
		t.Error("expected to find 'feature-2' multidev")
	}
}

func TestMultidevService_ListEmpty(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return only standard environments
		envs := []map[string]interface{}{
			{"id": "dev", "domain": "dev.site1.pantheonsite.io"},
			{"id": "test", "domain": "test.site1.pantheonsite.io"},
			{"id": "live", "domain": "live.site1.pantheonsite.io"},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(envs)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	multidevService := NewMultidevService(client)

	multidevs, err := multidevService.List(context.Background(), "site1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return empty list when no multidevs exist
	if len(multidevs) != 0 {
		t.Errorf("expected 0 multidev environments, got %d", len(multidevs))
	}
}
