package api

import (
	"context"
	"fmt"
)

// SessionTokenRefresher implements TokenRefresher using machine token authentication.
// It refreshes the session token by calling the authentication endpoint with the
// stored machine token.
type SessionTokenRefresher struct {
	// machineToken is the machine token used to obtain new session tokens
	machineToken string
	// client is the API client used to make the refresh request
	// Note: This creates a circular reference, but it's necessary for the refresh flow.
	// The doWithRetry function tracks tokenRefreshAttempted to prevent infinite loops.
	client *Client
	// onTokenRefreshed is an optional callback invoked when a token is successfully refreshed.
	// This allows the caller to persist the new session information.
	onTokenRefreshed func(session *SessionResponse) error
	// logger is an optional logger for debug output
	logger Logger
}

// SessionTokenRefresherOption is a function that configures a SessionTokenRefresher
type SessionTokenRefresherOption func(*SessionTokenRefresher)

// WithRefreshLogger sets the logger for the token refresher
func WithRefreshLogger(logger Logger) SessionTokenRefresherOption {
	return func(r *SessionTokenRefresher) {
		r.logger = logger
	}
}

// WithOnTokenRefreshed sets the callback invoked when a token is refreshed
func WithOnTokenRefreshed(callback func(session *SessionResponse) error) SessionTokenRefresherOption {
	return func(r *SessionTokenRefresher) {
		r.onTokenRefreshed = callback
	}
}

// NewSessionTokenRefresher creates a new SessionTokenRefresher.
// The machineToken is required and will be used to obtain new session tokens.
// The client is the API client to use for the refresh request.
func NewSessionTokenRefresher(machineToken string, client *Client, opts ...SessionTokenRefresherOption) *SessionTokenRefresher {
	r := &SessionTokenRefresher{
		machineToken: machineToken,
		client:       client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// SetMachineToken updates the machine token used for refreshing
func (r *SessionTokenRefresher) SetMachineToken(token string) {
	r.machineToken = token
}

// RefreshToken implements the TokenRefresher interface.
// It uses the stored machine token to obtain a new session token from the API.
func (r *SessionTokenRefresher) RefreshToken(ctx context.Context) (string, error) {
	if r.machineToken == "" {
		return "", fmt.Errorf("no machine token available for token refresh")
	}

	if r.logger != nil {
		r.logger.Debug("Refreshing session token using machine token")
	}

	// Create a temporary auth service to perform the login
	// Note: We use the existing client which will use the same base URL and settings.
	// The doWithRetry function in the client tracks tokenRefreshAttempted to prevent
	// infinite recursion if the auth endpoint also returns 401.
	authService := NewAuthService(r.client)
	session, err := authService.Login(ctx, r.machineToken)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	if r.logger != nil {
		r.logger.Debug("Session token refreshed successfully")
	}

	// Call the callback if set (to persist the new session)
	if r.onTokenRefreshed != nil {
		if err := r.onTokenRefreshed(session); err != nil {
			if r.logger != nil {
				r.logger.Warn("Failed to save refreshed session: %v", err)
			}
			// Continue anyway since we have the new token
		}
	}

	return session.Session, nil
}
