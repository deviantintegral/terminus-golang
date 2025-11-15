package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Session represents an authenticated session
type Session struct {
	SessionToken string `json:"session"`
	UserID       string `json:"user_id"`
	Email        string `json:"email,omitempty"`
	ExpiresAt    int64  `json:"expires_at"`
}

// IsExpired returns true if the session has expired
func (s *Session) IsExpired() bool {
	if s.ExpiresAt == 0 {
		return false
	}
	return time.Now().Unix() > s.ExpiresAt
}

// Store handles session and token persistence
type Store struct {
	sessionPath string
	tokensPath  string
}

// NewStore creates a new session store
func NewStore(cacheDir string) *Store {
	return &Store{
		sessionPath: filepath.Join(cacheDir, "session"),
		tokensPath:  filepath.Join(cacheDir, "tokens"),
	}
}

// SaveSession saves a session to disk
func (s *Store) SaveSession(session *Session) error {
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(s.sessionPath), 0700); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	// Write with secure permissions
	if err := os.WriteFile(s.sessionPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// LoadSession loads a session from disk
func (s *Store) LoadSession() (*Session, error) {
	data, err := os.ReadFile(s.sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Check if expired
	if session.IsExpired() {
		// Remove expired session
		_ = s.DeleteSession()
		return nil, nil
	}

	return &session, nil
}

// DeleteSession deletes the session file
func (s *Store) DeleteSession() error {
	if err := os.Remove(s.sessionPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete session file: %w", err)
	}
	return nil
}

// SaveToken saves a machine token to disk
func (s *Store) SaveToken(email, token string) error {
	// Ensure tokens directory exists
	if err := os.MkdirAll(s.tokensPath, 0700); err != nil {
		return fmt.Errorf("failed to create tokens directory: %w", err)
	}

	tokenPath := filepath.Join(s.tokensPath, sanitizeFilename(email))

	// Write with secure permissions
	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadToken loads a machine token from disk
func (s *Store) LoadToken(email string) (string, error) {
	tokenPath := filepath.Join(s.tokensPath, sanitizeFilename(email))

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	return string(data), nil
}

// DeleteToken deletes a machine token
func (s *Store) DeleteToken(email string) error {
	tokenPath := filepath.Join(s.tokensPath, sanitizeFilename(email))

	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	return nil
}

// ListTokens returns a list of saved token emails
func (s *Store) ListTokens() ([]string, error) {
	entries, err := os.ReadDir(s.tokensPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read tokens directory: %w", err)
	}

	var emails []string
	for _, entry := range entries {
		if !entry.IsDir() {
			emails = append(emails, entry.Name())
		}
	}

	return emails, nil
}

// sanitizeFilename sanitizes a filename to prevent path traversal
func sanitizeFilename(name string) string {
	// Replace path separators and other problematic characters
	safe := filepath.Base(name)
	safe = filepath.Clean(safe)
	return safe
}
