package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSessionTokenRefresher(t *testing.T) {
	client := NewClient()
	refresher := NewSessionTokenRefresher("machine-token", client)

	if refresher == nil {
		t.Fatal("expected refresher to be created")
	}

	if refresher.machineToken != "machine-token" {
		t.Errorf("expected machine token 'machine-token', got %s", refresher.machineToken)
	}

	if refresher.client != client {
		t.Error("expected client to be set")
	}
}

func TestSessionTokenRefresher_RefreshToken_Success(t *testing.T) {
	// Create a mock auth server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify it's the auth endpoint
		if r.URL.Path != "/authorize/machine-token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Verify method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		// Verify machine token in body
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.MachineToken != "test-machine-token" {
			t.Errorf("expected machine token 'test-machine-token', got %s", req.MachineToken)
		}

		// Return session response
		resp := SessionResponse{
			Session:   "new-session-token",
			UserID:    "user-123",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	var savedSession *SessionResponse
	refresher := NewSessionTokenRefresher(
		"test-machine-token",
		client,
		WithOnTokenRefreshed(func(session *SessionResponse) error {
			savedSession = session
			return nil
		}),
	)

	ctx := context.Background()
	token, err := refresher.RefreshToken(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token != "new-session-token" {
		t.Errorf("expected token 'new-session-token', got %s", token)
	}

	// Verify callback was called
	if savedSession == nil {
		t.Fatal("expected onTokenRefreshed callback to be called")
	}

	if savedSession.Session != "new-session-token" {
		t.Errorf("expected saved session token 'new-session-token', got %s", savedSession.Session)
	}
}

func TestSessionTokenRefresher_RefreshToken_NoMachineToken(t *testing.T) {
	client := NewClient()
	refresher := NewSessionTokenRefresher("", client)

	ctx := context.Background()
	_, err := refresher.RefreshToken(ctx)
	if err == nil {
		t.Fatal("expected error for empty machine token")
	}

	expectedErr := "no machine token available for token refresh"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestSessionTokenRefresher_RefreshToken_AuthFails(t *testing.T) {
	// Create a mock auth server that returns 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "invalid machine token"}`))
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	refresher := NewSessionTokenRefresher("invalid-token", client)

	ctx := context.Background()
	_, err := refresher.RefreshToken(ctx)
	if err == nil {
		t.Fatal("expected error for failed auth")
	}
}

func TestSessionTokenRefresher_SetMachineToken(t *testing.T) {
	client := NewClient()
	refresher := NewSessionTokenRefresher("original-token", client)

	if refresher.machineToken != "original-token" {
		t.Errorf("expected 'original-token', got %s", refresher.machineToken)
	}

	refresher.SetMachineToken("new-token")

	if refresher.machineToken != "new-token" {
		t.Errorf("expected 'new-token', got %s", refresher.machineToken)
	}
}

func TestSessionTokenRefresher_WithRefreshLogger(t *testing.T) {
	client := NewClient()
	logger := NewLogger(VerbosityDebug)

	refresher := NewSessionTokenRefresher(
		"token",
		client,
		WithRefreshLogger(logger),
	)

	if refresher.logger == nil {
		t.Error("expected logger to be set")
	}
}
