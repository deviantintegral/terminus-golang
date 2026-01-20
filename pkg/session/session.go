// Package session handles session and token storage for authentication.
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
	ExpiresAt    int64  `json:"expires_at"`
	Email        string `json:"email,omitempty"`
}

// IsExpired returns true if the session has expired
func (s *Session) IsExpired() bool {
	if s.ExpiresAt == 0 {
		return false
	}
	return time.Now().Unix() > s.ExpiresAt
}

// TokenRenewalBuffer is the duration before expiry when a token is considered
// to need renewal (5 minutes)
const TokenRenewalBuffer = 5 * time.Minute

// NeedsRenewal returns true if the session token will expire within the
// renewal buffer window (5 minutes) or has already expired
func (s *Session) NeedsRenewal() bool {
	if s.ExpiresAt == 0 {
		return false // No expiry time set, assume indefinite validity
	}
	renewalTime := time.Now().Add(TokenRenewalBuffer).Unix()
	return renewalTime > s.ExpiresAt
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
	if err := os.MkdirAll(filepath.Dir(s.sessionPath), 0o700); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	// Write with secure permissions
	if err := os.WriteFile(s.sessionPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// LoadSession loads a session from disk
// Note: This returns expired sessions as well, since the machine token
// can still be used for automatic token renewal. Callers should check
// IsExpired() if they need to verify the session is still valid.
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

	return &session, nil
}

// DeleteSession deletes the session file
func (s *Store) DeleteSession() error {
	if err := os.Remove(s.sessionPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete session file: %w", err)
	}
	return nil
}

// SaveToken saves a machine token to disk in PHP Terminus compatible JSON format
func (s *Store) SaveToken(email, token string) error {
	// Ensure tokens directory exists
	if err := os.MkdirAll(s.tokensPath, 0o700); err != nil {
		return fmt.Errorf("failed to create tokens directory: %w", err)
	}

	tokenPath := filepath.Join(s.tokensPath, sanitizeFilename(email))

	// Save in PHP Terminus JSON format for compatibility
	tokenData := phpTokenFormat{
		Token: token,
		Email: email,
		Date:  time.Now().Unix(),
	}

	data, err := json.MarshalIndent(tokenData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Write with secure permissions
	if err := os.WriteFile(tokenPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadToken loads a machine token from disk (returns raw file content)
func (s *Store) LoadToken(email string) (string, error) {
	tokenPath := filepath.Join(s.tokensPath, sanitizeFilename(email))

	data, err := os.ReadFile(tokenPath) //nolint:gosec // Token file path
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	return string(data), nil
}

// LoadMachineToken loads and extracts the raw machine token for an email.
// This handles both PHP Terminus JSON format and raw token strings.
func (s *Store) LoadMachineToken(email string) (string, error) {
	tokenData, err := s.LoadToken(email)
	if err != nil {
		return "", err
	}
	if tokenData == "" {
		return "", nil
	}
	return ExtractRawToken(tokenData), nil
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

// phpTokenFormat represents the JSON format used by PHP Terminus for stored tokens
type phpTokenFormat struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Date  int64  `json:"date"`
}

// ExtractRawToken extracts the raw machine token value from a token string.
// PHP Terminus stores tokens as JSON with token, email, and date fields.
// This function handles both formats: raw token strings and PHP-style JSON.
func ExtractRawToken(tokenData string) string {
	// Try to parse as PHP Terminus JSON format
	var phpToken phpTokenFormat
	if err := json.Unmarshal([]byte(tokenData), &phpToken); err == nil {
		// Successfully parsed as PHP format, return the token (even if empty)
		return phpToken.Token
	}
	// Return as-is if not in PHP format
	return tokenData
}
