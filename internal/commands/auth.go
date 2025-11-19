// Package commands implements all CLI commands for terminus.
package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/pantheon-systems/terminus-go/pkg/session"
	"github.com/spf13/cobra"
)

var authLoginCmd = &cobra.Command{
	Use:   "auth:login",
	Short: "Log in to Pantheon",
	Long: `Authenticate with Pantheon using a machine token.

Usage examples:
  auth:login --machine-token=<machine_token>  Logs in with the provided machine token
  auth:login                                   Logs in with a previously saved machine token
  auth:login --email=<email>                   Logs in with a saved token belonging to <email>`,
	RunE: runAuthLogin,
}

var authLogoutCmd = &cobra.Command{
	Use:   "auth:logout",
	Short: "Log out of Pantheon",
	Long:  "Remove stored authentication credentials",
	RunE:  runAuthLogout,
}

var authWhoamiCmd = &cobra.Command{
	Use:   "auth:whoami",
	Short: "Show current user",
	Long:  "Display information about the currently authenticated user",
	RunE:  runAuthWhoami,
}

var (
	machineTokenFlag string
	emailFlag        string
)

func init() {
	// Add auth commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(authLoginCmd)
	rootCmd.AddCommand(authLogoutCmd)
	rootCmd.AddCommand(authWhoamiCmd)

	authLoginCmd.Flags().StringVar(&machineTokenFlag, "machine-token", "", "Machine token for authentication")
	authLoginCmd.Flags().StringVar(&emailFlag, "email", "", "Email address (for token lookup/storage)")
}

func runAuthLogin(_ *cobra.Command, _ []string) error {
	// Determine the machine token to use
	token := machineTokenFlag
	email := emailFlag

	if token == "" {
		// No token provided, try to load from saved tokens
		var err error
		token, email, err = resolveToken(email)
		if err != nil {
			return err
		}
	}

	// Extract raw token value (handles PHP Terminus JSON format)
	token = session.ExtractRawToken(token)

	// Create auth service
	authService := api.NewAuthService(cliContext.APIClient)

	// Authenticate
	printMessage("Logging in...")
	sess, err := authService.Login(getContext(), token)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Save session
	sessionData := &session.Session{
		SessionToken: sess.Session,
		UserID:       sess.UserID,
		ExpiresAt:    sess.ExpiresAt,
	}

	if err := cliContext.SessionStore.SaveSession(sessionData); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	// Save machine token for future use if email was provided
	if email != "" {
		if err := cliContext.SessionStore.SaveToken(email, token); err != nil {
			printError("Warning: failed to save machine token: %v", err)
		}
	}

	// Update API client with new session token
	cliContext.APIClient.SetToken(sess.Session)

	printMessage("Login successful!")

	// Get and display user info
	user, err := authService.Whoami(getContext(), sess.UserID)
	if err != nil {
		printError("Warning: failed to get user info: %v", err)
		return nil
	}

	printMessage("Logged in as: %s %s (%s)", user.FirstName, user.LastName, user.Email)

	return nil
}

// resolveToken loads a machine token from saved tokens
func resolveToken(email string) (token, resolvedEmail string, err error) {
	if email != "" {
		// Load token for specific email
		token, err = cliContext.SessionStore.LoadToken(email)
		if err != nil {
			return "", "", fmt.Errorf("failed to load token: %w", err)
		}
		if token == "" {
			return "", "", fmt.Errorf("no saved token found for %s", email)
		}
		return token, email, nil
	}

	// No email provided, list all saved tokens
	var emails []string
	emails, err = cliContext.SessionStore.ListTokens()
	if err != nil {
		return "", "", fmt.Errorf("failed to list tokens: %w", err)
	}

	if len(emails) == 0 {
		return "", "", fmt.Errorf("no saved machine tokens found. Please provide --machine-token")
	}

	if len(emails) > 1 {
		return "", "", fmt.Errorf("multiple saved tokens found. Please specify --email with one of: %v", emails)
	}

	// Exactly one saved token
	resolvedEmail = emails[0]
	token, err = cliContext.SessionStore.LoadToken(resolvedEmail)
	if err != nil {
		return "", "", fmt.Errorf("failed to load token for %s: %w", resolvedEmail, err)
	}
	if token == "" {
		return "", "", fmt.Errorf("saved token for %s is empty", resolvedEmail)
	}

	printMessage("Using saved token for %s", resolvedEmail)
	return token, resolvedEmail, nil
}

func runAuthLogout(_ *cobra.Command, _ []string) error {
	// Delete session
	if err := cliContext.SessionStore.DeleteSession(); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	printMessage("Logged out successfully")

	return nil
}

func runAuthWhoami(_ *cobra.Command, _ []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	// Load session to get user ID
	sess, err := cliContext.SessionStore.LoadSession()
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}
	if sess == nil || sess.UserID == "" {
		return fmt.Errorf("no user ID in session")
	}

	// Create auth service
	authService := api.NewAuthService(cliContext.APIClient)

	// Get user info
	user, err := authService.Whoami(getContext(), sess.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	return printOutput(user)
}
