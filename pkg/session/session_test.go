package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSessionIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt int64
		expected  bool
	}{
		{
			name:      "not expired",
			expiresAt: time.Now().Add(1 * time.Hour).Unix(),
			expected:  false,
		},
		{
			name:      "expired",
			expiresAt: time.Now().Add(-1 * time.Hour).Unix(),
			expected:  true,
		},
		{
			name:      "no expiry",
			expiresAt: 0,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sess := &Session{
				SessionToken: "test-token",
				UserID:       "test-user",
				ExpiresAt:    tt.expiresAt,
			}

			if sess.IsExpired() != tt.expected {
				t.Errorf("expected IsExpired() to be %v, got %v", tt.expected, sess.IsExpired())
			}
		})
	}
}

func TestStoreSaveAndLoadSession(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	sess := &Session{
		SessionToken: "test-token",
		UserID:       "test-user-id",
		ExpiresAt:    time.Now().Add(1 * time.Hour).Unix(),
	}

	// Test save
	err := store.SaveSession(sess)
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// Test load
	loaded, err := store.LoadSession()
	if err != nil {
		t.Fatalf("failed to load session: %v", err)
	}

	if loaded == nil {
		t.Fatal("expected session to be loaded")
	}

	if loaded.SessionToken != sess.SessionToken {
		t.Errorf("expected token %s, got %s", sess.SessionToken, loaded.SessionToken)
	}

	if loaded.UserID != sess.UserID {
		t.Errorf("expected user ID %s, got %s", sess.UserID, loaded.UserID)
	}
}

func TestStoreExpiredSession(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	sess := &Session{
		SessionToken: "test-token",
		UserID:       "test-user-id",
		ExpiresAt:    time.Now().Add(-1 * time.Hour).Unix(),
	}

	// Save expired session
	err := store.SaveSession(sess)
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// Try to load expired session
	loaded, err := store.LoadSession()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expired session should return nil
	if loaded != nil {
		t.Error("expected expired session to return nil")
	}

	// Verify session file was deleted
	_, err = os.Stat(store.sessionPath)
	if !os.IsNotExist(err) {
		t.Error("expected expired session file to be deleted")
	}
}

func TestStoreDeleteSession(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	sess := &Session{
		SessionToken: "test-token",
		UserID:       "test-user-id",
		ExpiresAt:    time.Now().Add(1 * time.Hour).Unix(),
	}

	// Save session
	err := store.SaveSession(sess)
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// Delete session
	err = store.DeleteSession()
	if err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}

	// Verify session file is gone
	_, err = os.Stat(store.sessionPath)
	if !os.IsNotExist(err) {
		t.Error("expected session file to be deleted")
	}
}

func TestStoreSaveAndLoadToken(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	email := "test@example.com"
	token := "test-machine-token"

	// Test save
	err := store.SaveToken(email, token)
	if err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	// Test load
	loaded, err := store.LoadToken(email)
	if err != nil {
		t.Fatalf("failed to load token: %v", err)
	}

	if loaded != token {
		t.Errorf("expected token %s, got %s", token, loaded)
	}
}

func TestStoreListTokens(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	emails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}

	// Save multiple tokens
	for _, email := range emails {
		err := store.SaveToken(email, "token-"+email)
		if err != nil {
			t.Fatalf("failed to save token for %s: %v", email, err)
		}
	}

	// List tokens
	list, err := store.ListTokens()
	if err != nil {
		t.Fatalf("failed to list tokens: %v", err)
	}

	if len(list) != len(emails) {
		t.Errorf("expected %d tokens, got %d", len(emails), len(list))
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test@example.com", "test@example.com"},
		{"../../../etc/passwd", "passwd"},
		{"/absolute/path", "path"},
		{"normal-filename.txt", "normal-filename.txt"},
	}

	for _, tt := range tests {
		result := sanitizeFilename(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeFilename(%s) = %s, expected %s", tt.input, result, tt.expected)
		}

		// Ensure result doesn't contain path separators
		if filepath.Dir(result) != "." {
			t.Errorf("sanitized filename %s still contains path separators", result)
		}
	}
}

func TestExtractRawToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "raw token string",
			input:    "abcd1234-raw-token-value",
			expected: "abcd1234-raw-token-value",
		},
		{
			name:     "PHP Terminus JSON format",
			input:    `{"token":"actual-token-value","email":"user@example.com","date":1763142241}`,
			expected: "actual-token-value",
		},
		{
			name:     "PHP format with different field order",
			input:    `{"email":"user@example.com","date":1763142241,"token":"my-token"}`,
			expected: "my-token",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "invalid JSON",
			input:    `{"token": incomplete`,
			expected: `{"token": incomplete`,
		},
		{
			name:     "JSON without token field",
			input:    `{"email":"user@example.com","date":1763142241}`,
			expected: `{"email":"user@example.com","date":1763142241}`,
		},
		{
			name:     "JSON with empty token field",
			input:    `{"token":"","email":"user@example.com"}`,
			expected: `{"token":"","email":"user@example.com"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractRawToken(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractRawToken(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}
