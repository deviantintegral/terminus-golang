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
	Long:  "Authenticate with Pantheon using a machine token",
	RunE:  runAuthLogin,
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
	saveTokenFlag    bool
)

func init() {
	// Add auth commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(authLoginCmd)
	rootCmd.AddCommand(authLogoutCmd)
	rootCmd.AddCommand(authWhoamiCmd)

	authLoginCmd.Flags().StringVar(&machineTokenFlag, "machine-token", "", "Machine token for authentication")
	authLoginCmd.Flags().StringVar(&emailFlag, "email", "", "Email address (for token storage)")
	authLoginCmd.Flags().BoolVar(&saveTokenFlag, "save-token", true, "Save machine token for future use")
	_ = authLoginCmd.MarkFlagRequired("machine-token")
}

func runAuthLogin(_ *cobra.Command, _ []string) error {
	// Create auth service
	authService := api.NewAuthService(cliContext.APIClient)

	// Authenticate
	printMessage("Logging in...")
	sess, err := authService.Login(getContext(), machineTokenFlag, emailFlag)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Save session
	sessionData := &session.Session{
		SessionToken: sess.Session,
		UserID:       sess.UserID,
		Email:        sess.Email,
		ExpiresAt:    sess.ExpiresAt,
	}

	if err := cliContext.SessionStore.SaveSession(sessionData); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	// Save machine token if requested
	if saveTokenFlag && emailFlag != "" {
		if err := cliContext.SessionStore.SaveToken(emailFlag, machineTokenFlag); err != nil {
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
