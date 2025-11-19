package commands

import (
	"fmt"
	"testing"

	"github.com/pantheon-systems/terminus-go/pkg/session"
)

func TestAuthLoginCmdStructure(t *testing.T) {
	if authLoginCmd.Use != "auth:login" {
		t.Errorf("expected authLoginCmd.Use to be 'auth:login', got '%s'", authLoginCmd.Use)
	}

	if authLoginCmd.Short == "" {
		t.Error("authLoginCmd.Short should not be empty")
	}

	if authLoginCmd.Long == "" {
		t.Error("authLoginCmd.Long should not be empty")
	}

	// Check flags
	machineTokenFlag := authLoginCmd.Flags().Lookup("machine-token")
	if machineTokenFlag == nil {
		t.Error("authLoginCmd should have a 'machine-token' flag")
	}

	emailFlag := authLoginCmd.Flags().Lookup("email")
	if emailFlag == nil {
		t.Error("authLoginCmd should have an 'email' flag")
	}
}

func TestAuthLogoutCmdStructure(t *testing.T) {
	if authLogoutCmd.Use != "auth:logout" {
		t.Errorf("expected authLogoutCmd.Use to be 'auth:logout', got '%s'", authLogoutCmd.Use)
	}

	if authLogoutCmd.Short == "" {
		t.Error("authLogoutCmd.Short should not be empty")
	}

	if authLogoutCmd.Long == "" {
		t.Error("authLogoutCmd.Long should not be empty")
	}
}

func TestAuthWhoamiCmdStructure(t *testing.T) {
	if authWhoamiCmd.Use != "auth:whoami" {
		t.Errorf("expected authWhoamiCmd.Use to be 'auth:whoami', got '%s'", authWhoamiCmd.Use)
	}

	if authWhoamiCmd.Short == "" {
		t.Error("authWhoamiCmd.Short should not be empty")
	}

	if authWhoamiCmd.Long == "" {
		t.Error("authWhoamiCmd.Long should not be empty")
	}
}

func TestAuthCommands(t *testing.T) {
	expectedCommands := []string{"auth:login", "auth:logout", "auth:whoami"}

	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected command '%s' not found in rootCmd", expected)
		}
	}
}

func TestResolveToken_WithEmail(t *testing.T) {
	// Create a temporary session store
	tmpDir := t.TempDir()
	store := session.NewStore(tmpDir)

	// Save a token for a specific email
	email := "test@example.com"
	token := "test-token-123"
	if err := store.SaveToken(email, token); err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	// Set up CLI context
	oldContext := cliContext
	defer func() { cliContext = oldContext }()
	cliContext = &CLIContext{
		SessionStore: store,
	}

	// Test resolving token with specific email
	resolvedToken, resolvedEmail, err := resolveToken(email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedToken != token {
		t.Errorf("expected token '%s', got '%s'", token, resolvedToken)
	}

	if resolvedEmail != email {
		t.Errorf("expected email '%s', got '%s'", email, resolvedEmail)
	}
}

func TestResolveToken_WithEmail_NotFound(t *testing.T) {
	// Create a temporary session store
	tmpDir := t.TempDir()
	store := session.NewStore(tmpDir)

	// Set up CLI context
	oldContext := cliContext
	defer func() { cliContext = oldContext }()
	cliContext = &CLIContext{
		SessionStore: store,
	}

	// Test resolving token for non-existent email
	_, _, err := resolveToken("nonexistent@example.com")
	if err == nil {
		t.Fatal("expected error for non-existent email")
	}

	expectedErr := "no saved token found for nonexistent@example.com"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestResolveToken_NoEmail_SingleToken(t *testing.T) {
	// Create a temporary session store
	tmpDir := t.TempDir()
	store := session.NewStore(tmpDir)

	// Save a single token
	email := "test@example.com"
	token := "test-token-123"
	if err := store.SaveToken(email, token); err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	// Set up CLI context
	oldContext := cliContext
	defer func() { cliContext = oldContext }()
	cliContext = &CLIContext{
		SessionStore: store,
	}

	// Test resolving token with no email (should use the single saved token)
	resolvedToken, resolvedEmail, err := resolveToken("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedToken != token {
		t.Errorf("expected token '%s', got '%s'", token, resolvedToken)
	}

	if resolvedEmail != email {
		t.Errorf("expected email '%s', got '%s'", email, resolvedEmail)
	}
}

func TestResolveToken_NoEmail_NoTokens(t *testing.T) {
	// Create a temporary session store
	tmpDir := t.TempDir()
	store := session.NewStore(tmpDir)

	// Set up CLI context
	oldContext := cliContext
	defer func() { cliContext = oldContext }()
	cliContext = &CLIContext{
		SessionStore: store,
	}

	// Test resolving token with no email and no saved tokens
	_, _, err := resolveToken("")
	if err == nil {
		t.Fatal("expected error for no saved tokens")
	}

	expectedErr := "no saved machine tokens found. Please provide --machine-token"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestResolveToken_NoEmail_MultipleTokens(t *testing.T) {
	// Create a temporary session store
	tmpDir := t.TempDir()
	store := session.NewStore(tmpDir)

	// Save multiple tokens
	tokens := map[string]string{
		"test1@example.com": "token1",
		"test2@example.com": "token2",
		"test3@example.com": "token3",
	}
	for email, token := range tokens {
		if err := store.SaveToken(email, token); err != nil {
			t.Fatalf("failed to save token: %v", err)
		}
	}

	// Set up CLI context
	oldContext := cliContext
	defer func() { cliContext = oldContext }()
	cliContext = &CLIContext{
		SessionStore: store,
	}

	// Test resolving token with no email and multiple saved tokens
	_, _, err := resolveToken("")
	if err == nil {
		t.Fatal("expected error for multiple saved tokens")
	}

	// Error should mention multiple tokens
	if err.Error()[:38] != "multiple saved tokens found. Please sp" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestResolveToken_LoadTokenError(t *testing.T) {
	// Create a temporary session store
	tmpDir := t.TempDir()
	store := session.NewStore(tmpDir)

	// Save a token
	email := "test@example.com"
	token := "test-token-123"
	if err := store.SaveToken(email, token); err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	// Set up CLI context with different store to simulate error
	badStore := session.NewStore("/nonexistent/path/that/does/not/exist")
	oldContext := cliContext
	defer func() { cliContext = oldContext }()
	cliContext = &CLIContext{
		SessionStore: badStore,
	}

	// Test resolving token with email (should fail to load)
	_, _, err := resolveToken(email)
	if err == nil {
		t.Fatal("expected error when loading token fails")
	}
}

func TestResolveToken_EmptyToken(t *testing.T) {
	// Create a temporary session store
	tmpDir := t.TempDir()
	store := session.NewStore(tmpDir)

	// Save an empty token (shouldn't happen in practice, but testing edge case)
	email := "test@example.com"
	if err := store.SaveToken(email, ""); err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	// Set up CLI context
	oldContext := cliContext
	defer func() { cliContext = oldContext }()
	cliContext = &CLIContext{
		SessionStore: store,
	}

	// Test resolving token with email that has empty token
	_, _, err := resolveToken(email)
	if err == nil {
		t.Fatal("expected error for empty token")
	}

	expectedErr := fmt.Sprintf("no saved token found for %s", email)
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestResolveToken_SingleEmptyToken(t *testing.T) {
	// Create a temporary session store
	tmpDir := t.TempDir()
	store := session.NewStore(tmpDir)

	// Save an empty token
	email := "test@example.com"
	if err := store.SaveToken(email, ""); err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	// Set up CLI context
	oldContext := cliContext
	defer func() { cliContext = oldContext }()
	cliContext = &CLIContext{
		SessionStore: store,
	}

	// Test resolving token with no email when single token is empty
	_, _, err := resolveToken("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}

	expectedErr := fmt.Sprintf("saved token for %s is empty", email)
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestAuthLoginCmd_RunEPresent(t *testing.T) {
	// Verify that the RunE function is set
	if authLoginCmd.RunE == nil {
		t.Error("authLoginCmd.RunE should not be nil")
	}
}

func TestAuthLogoutCmd_RunEPresent(t *testing.T) {
	// Verify that the RunE function is set
	if authLogoutCmd.RunE == nil {
		t.Error("authLogoutCmd.RunE should not be nil")
	}
}

func TestAuthWhoamiCmd_RunEPresent(t *testing.T) {
	// Verify that the RunE function is set
	if authWhoamiCmd.RunE == nil {
		t.Error("authWhoamiCmd.RunE should not be nil")
	}
}

func TestAuthLoginCmd_Flags(t *testing.T) {
	// Test machine-token flag
	flag := authLoginCmd.Flags().Lookup("machine-token")
	if flag == nil {
		t.Fatal("machine-token flag not found")
	}

	if flag.DefValue != "" {
		t.Errorf("expected machine-token default to be empty, got '%s'", flag.DefValue)
	}

	// Test email flag
	emailFlag := authLoginCmd.Flags().Lookup("email")
	if emailFlag == nil {
		t.Fatal("email flag not found")
	}

	if emailFlag.DefValue != "" {
		t.Errorf("expected email default to be empty, got '%s'", emailFlag.DefValue)
	}
}
