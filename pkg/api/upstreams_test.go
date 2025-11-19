package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpstreamsService_List(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/test-user/upstreams" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		upstreams := []map[string]interface{}{
			{
				"id":              "upstream1",
				"label":           "Drupal 10",
				"machine_name":    "drupal10",
				"type":            "core",
				"framework":       "drupal",
				"organization_id": "org1",
			},
			{
				"id":              "upstream2",
				"label":           "WordPress",
				"machine_name":    "wordpress",
				"type":            "core",
				"framework":       "wordpress",
				"organization_id": "org2",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(upstreams)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	upstreamsService := NewUpstreamsService(client)

	upstreams, err := upstreamsService.List(context.Background(), "test-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(upstreams) != 2 {
		t.Errorf("expected 2 upstreams, got %d", len(upstreams))
	}

	if upstreams[0].ID != "upstream1" {
		t.Errorf("expected upstream ID 'upstream1', got '%s'", upstreams[0].ID)
	}

	if upstreams[0].Label != "Drupal 10" {
		t.Errorf("expected label 'Drupal 10', got '%s'", upstreams[0].Label)
	}

	if upstreams[0].Framework != "drupal" {
		t.Errorf("expected framework 'drupal', got '%s'", upstreams[0].Framework)
	}
}

func TestUpstreamsService_Get(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/upstreams/upstream1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		upstream := map[string]interface{}{
			"id":              "upstream1",
			"label":           "Drupal 10",
			"machine_name":    "drupal10",
			"type":            "core",
			"framework":       "drupal",
			"organization_id": "org1",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(upstream)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	upstreamsService := NewUpstreamsService(client)

	upstream, err := upstreamsService.Get(context.Background(), "upstream1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if upstream.ID != "upstream1" {
		t.Errorf("expected upstream ID 'upstream1', got '%s'", upstream.ID)
	}

	if upstream.Label != "Drupal 10" {
		t.Errorf("expected label 'Drupal 10', got '%s'", upstream.Label)
	}
}

func TestUpstreamsService_ListUpdates(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sites/site1/environments/dev/code-upstream-updates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		updates := []map[string]interface{}{
			{
				"hash":     "abc123",
				"datetime": "2024-01-15T10:30:00Z",
				"message":  "Update module",
				"author":   "Developer",
			},
			{
				"hash":     "def456",
				"datetime": "2024-01-14T09:00:00Z",
				"message":  "Security patch",
				"author":   "Security Team",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(updates)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	upstreamsService := NewUpstreamsService(client)

	updates, err := upstreamsService.ListUpdates(context.Background(), "site1", "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(updates) != 2 {
		t.Errorf("expected 2 updates, got %d", len(updates))
	}

	if updates[0].Hash != "abc123" {
		t.Errorf("expected hash 'abc123', got '%s'", updates[0].Hash)
	}

	if updates[0].Message != "Update module" {
		t.Errorf("expected message 'Update module', got '%s'", updates[0].Message)
	}
}
