package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUsersService_ListMachineTokens(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/test-user/machine_tokens" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		tokens := []map[string]interface{}{
			{"id": "token1", "device_name": "Device 1"},
			{"id": "token2", "device_name": "Device 2"},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tokens)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	usersService := NewUsersService(client)

	tokens, err := usersService.ListMachineTokens(context.Background(), "test-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}

	if tokens[0].ID != "token1" {
		t.Errorf("expected token ID 'token1', got '%s'", tokens[0].ID)
	}

	if tokens[0].DeviceName != "Device 1" {
		t.Errorf("expected device name 'Device 1', got '%s'", tokens[0].DeviceName)
	}
}

func TestUsersService_ListSSHKeys(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/test-user/keys" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		keys := []map[string]interface{}{
			{"id": "key1", "hex": "aa:bb:cc", "key": "ssh-rsa AAAA..."},
			{"id": "key2", "hex": "dd:ee:ff", "key": "ssh-rsa BBBB..."},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(keys)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	usersService := NewUsersService(client)

	keys, err := usersService.ListSSHKeys(context.Background(), "test-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}

	if keys[0].ID != "key1" {
		t.Errorf("expected key ID 'key1', got '%s'", keys[0].ID)
	}

	if keys[0].Hex != "aa:bb:cc" {
		t.Errorf("expected hex 'aa:bb:cc', got '%s'", keys[0].Hex)
	}
}

func TestUsersService_ListPaymentMethods(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/test-user/instruments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		methods := []map[string]interface{}{
			{"id": "pm1", "label": "Visa ending in 1234"},
			{"id": "pm2", "label": "Mastercard ending in 5678"},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(methods)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	usersService := NewUsersService(client)

	methods, err := usersService.ListPaymentMethods(context.Background(), "test-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(methods) != 2 {
		t.Errorf("expected 2 methods, got %d", len(methods))
	}

	if methods[0].ID != "pm1" {
		t.Errorf("expected method ID 'pm1', got '%s'", methods[0].ID)
	}

	if methods[0].Label != "Visa ending in 1234" {
		t.Errorf("expected label 'Visa ending in 1234', got '%s'", methods[0].Label)
	}
}
