package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("expected client to be created")
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL to be %s, got %s", DefaultBaseURL, client.baseURL)
	}
}

func TestClientWithOptions(t *testing.T) {
	customBaseURL := "https://test.example.com"
	customToken := "test-token"

	client := NewClient(
		WithBaseURL(customBaseURL),
		WithToken(customToken),
	)

	if client.baseURL != customBaseURL {
		t.Errorf("expected baseURL to be %s, got %s", customBaseURL, client.baseURL)
	}

	if client.token != customToken {
		t.Errorf("expected token to be %s, got %s", customToken, client.token)
	}
}

func TestClientRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept header to be application/json, got %s", r.Header.Get("Accept"))
		}

		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization header to be 'Bearer test-token', got %s", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "test"}`))
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	ctx := context.Background()
	resp, err := client.Get(ctx, "/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestClientRequestErrorResponse(t *testing.T) {
	// Create a test server that returns a 404 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	ctx := context.Background()
	resp, err := client.Get(ctx, "/test")
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err == nil {
		t.Fatal("expected error for 404 response")
	}

	// Verify it's an API error
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", apiErr.StatusCode)
	}

	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
}

func TestClientRequestBadRequest(t *testing.T) {
	// Create a test server that returns a 400 error (no retries)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	ctx := context.Background()
	resp, err := client.Get(ctx, "/test")
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err == nil {
		t.Fatal("expected error for 400 response")
	}

	// Verify it's an API error with correct status
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code 400, got %d", apiErr.StatusCode)
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{400, false},
		{404, false},
		{429, true},
		{500, true},
		{502, true},
		{503, true},
	}

	for _, tt := range tests {
		result := shouldRetry(tt.statusCode)
		if result != tt.expected {
			t.Errorf("shouldRetry(%d) = %v, expected %v", tt.statusCode, result, tt.expected)
		}
	}
}

func TestError(t *testing.T) {
	err := &Error{
		StatusCode: 404,
		Message:    "Not Found",
	}

	expected := "API error 404: Not Found"
	if err.Error() != expected {
		t.Errorf("expected error message %s, got %s", expected, err.Error())
	}

	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true for 404 error")
	}

	if IsConflict(err) {
		t.Error("expected IsConflict to return false for 404 error")
	}
}
