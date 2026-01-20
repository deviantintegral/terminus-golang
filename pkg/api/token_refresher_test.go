package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSessionTokenRefresher(t *testing.T) {
	client := NewClient()
	getMachineToken := func() (string, error) {
		return "machine-token", nil
	}
	refresher := NewSessionTokenRefresher(getMachineToken, client)

	if refresher == nil {
		t.Fatal("expected refresher to be created")
	}

	if refresher.getMachineToken == nil {
		t.Error("expected getMachineToken callback to be set")
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
	getMachineToken := func() (string, error) {
		return "test-machine-token", nil
	}
	refresher := NewSessionTokenRefresher(
		getMachineToken,
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
	getMachineToken := func() (string, error) {
		return "", nil
	}
	refresher := NewSessionTokenRefresher(getMachineToken, client)

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

func TestSessionTokenRefresher_RefreshToken_NoProvider(t *testing.T) {
	client := NewClient()
	refresher := NewSessionTokenRefresher(nil, client)

	ctx := context.Background()
	_, err := refresher.RefreshToken(ctx)
	if err == nil {
		t.Fatal("expected error for nil machine token provider")
	}

	expectedErr := "no machine token provider configured"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestSessionTokenRefresher_RefreshToken_ProviderError(t *testing.T) {
	client := NewClient()
	getMachineToken := func() (string, error) {
		return "", errors.New("token file not found")
	}
	refresher := NewSessionTokenRefresher(getMachineToken, client)

	ctx := context.Background()
	_, err := refresher.RefreshToken(ctx)
	if err == nil {
		t.Fatal("expected error from provider")
	}

	if err.Error() != "failed to get machine token: token file not found" {
		t.Errorf("unexpected error: %s", err.Error())
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

	getMachineToken := func() (string, error) {
		return "invalid-token", nil
	}
	refresher := NewSessionTokenRefresher(getMachineToken, client)

	ctx := context.Background()
	_, err := refresher.RefreshToken(ctx)
	if err == nil {
		t.Fatal("expected error for failed auth")
	}
}

func TestSessionTokenRefresher_WithRefreshLogger(t *testing.T) {
	client := NewClient()
	logger := NewLogger(VerbosityDebug)

	getMachineToken := func() (string, error) {
		return "token", nil
	}
	refresher := NewSessionTokenRefresher(
		getMachineToken,
		client,
		WithRefreshLogger(logger),
	)

	if refresher.logger == nil {
		t.Error("expected logger to be set")
	}
}
