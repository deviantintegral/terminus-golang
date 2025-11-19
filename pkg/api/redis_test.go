package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRedisService(t *testing.T) {
	client := NewClient()
	service := NewRedisService(client)
	if service == nil {
		t.Fatal("expected service to be created")
	}
	if service.client != client {
		t.Error("expected service client to match provided client")
	}
}

func TestRedisService_Enable(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		expectedPath := "/sites/test-site/workflows"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody["type"] != "enable_addon" {
			t.Errorf("expected workflow type 'enable_addon', got %v", reqBody["type"])
		}

		params, ok := reqBody["params"].(map[string]interface{})
		if !ok {
			t.Fatal("expected params to be a map")
		}

		if params["addon"] != "cacheserver" {
			t.Errorf("expected addon 'cacheserver', got %v", params["addon"])
		}

		// Return successful workflow response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"id":          "workflow-123",
			"type":        "enable_addon",
			"description": "Enable Redis",
			"site_id":     "test-site",
			"result":      "",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewRedisService(client)

	ctx := context.Background()
	workflow, err := service.Enable(ctx, "test-site")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.ID != "workflow-123" {
		t.Errorf("expected workflow ID 'workflow-123', got %s", workflow.ID)
	}

	if workflow.Type != "enable_addon" {
		t.Errorf("expected workflow type 'enable_addon', got %s", workflow.Type)
	}
}

func TestRedisService_Disable(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		expectedPath := "/sites/test-site/workflows"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody["type"] != "disable_addon" {
			t.Errorf("expected workflow type 'disable_addon', got %v", reqBody["type"])
		}

		params, ok := reqBody["params"].(map[string]interface{})
		if !ok {
			t.Fatal("expected params to be a map")
		}

		if params["addon"] != "cacheserver" {
			t.Errorf("expected addon 'cacheserver', got %v", params["addon"])
		}

		// Return successful workflow response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"id":          "workflow-456",
			"type":        "disable_addon",
			"description": "Disable Redis",
			"site_id":     "test-site",
			"result":      "",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewRedisService(client)

	ctx := context.Background()
	workflow, err := service.Disable(ctx, "test-site")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.ID != "workflow-456" {
		t.Errorf("expected workflow ID 'workflow-456', got %s", workflow.ID)
	}

	if workflow.Type != "disable_addon" {
		t.Errorf("expected workflow type 'disable_addon', got %s", workflow.Type)
	}
}

func TestRedisService_Enable_Error(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Site not found"}`))
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewRedisService(client)

	ctx := context.Background()
	_, err := service.Enable(ctx, "nonexistent-site")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRedisService_Disable_Error(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Site not found"}`))
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewRedisService(client)

	ctx := context.Background()
	_, err := service.Disable(ctx, "nonexistent-site")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
