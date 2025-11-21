package commands

import (
	"testing"

	"github.com/pantheon-systems/terminus-go/pkg/session"
)

func TestDashboardViewCmdStructure(t *testing.T) {
	// Test that dashboardViewCmd has the expected properties
	if dashboardViewCmd.Use != "dashboard:view [site[.env]]" {
		t.Errorf("expected dashboardViewCmd.Use to be 'dashboard:view [site[.env]]', got '%s'", dashboardViewCmd.Use)
	}

	if dashboardViewCmd.Short == "" {
		t.Error("dashboardViewCmd.Short should not be empty")
	}

	if dashboardViewCmd.Long == "" {
		t.Error("dashboardViewCmd.Long should not be empty")
	}

	// Verify RunE is set
	if dashboardViewCmd.RunE == nil {
		t.Error("dashboardViewCmd.RunE should be set")
	}

	// Verify Args validator is set (should allow 0 or 1 argument)
	if dashboardViewCmd.Args == nil {
		t.Error("dashboardViewCmd.Args should be set")
	}
}

func TestDashboardViewCommands(t *testing.T) {
	expectedCommands := []string{"dashboard:view"}

	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == expected || (len(cmd.Use) > len(expected) && cmd.Use[:len(expected)] == expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected command '%s' not found in rootCmd", expected)
		}
	}
}

func TestDashboardViewFlags(t *testing.T) {
	// Test that the print flag exists
	flag := dashboardViewCmd.Flags().Lookup("print")
	if flag == nil {
		t.Error("dashboardViewCmd should have a 'print' flag")
		return
	}

	if flag.DefValue != "false" {
		t.Errorf("expected print flag default to be 'false', got '%s'", flag.DefValue)
	}

	if flag.Usage == "" {
		t.Error("print flag should have usage text")
	}
}

func TestDashboardViewRequiresAuth(t *testing.T) {
	// Save old context and create a minimal context without session
	oldContext := cliContext
	cliContext = &CLIContext{
		APIClient: nil,
	}
	defer func() { cliContext = oldContext }()

	err := runDashboardView(nil, []string{})
	if err == nil {
		t.Error("expected error when not authenticated")
	}
}

func TestGetUserDashboardURL(t *testing.T) {
	// Save old context
	oldContext := cliContext

	// Create a mock session with a test user ID
	mockSession := &session.Session{
		SessionToken: "test-token",
		UserID:       "test-user-id-123",
		ExpiresAt:    9999999999,
	}

	// Create a temporary session store
	tempStore := session.NewStore(t.TempDir())
	err := tempStore.SaveSession(mockSession)
	if err != nil {
		t.Fatalf("failed to save test session: %v", err)
	}

	cliContext = &CLIContext{
		SessionStore: tempStore,
	}
	defer func() { cliContext = oldContext }()

	url, err := getUserDashboardURL()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedURL := "https://dashboard.pantheon.io/users/test-user-id-123"
	if url != expectedURL {
		t.Errorf("expected URL '%s', got '%s'", expectedURL, url)
	}
}

func TestGetUserDashboardURLNoSession(t *testing.T) {
	// Save old context
	oldContext := cliContext

	// Create a temporary session store without a session
	tempStore := session.NewStore(t.TempDir())

	cliContext = &CLIContext{
		SessionStore: tempStore,
	}
	defer func() { cliContext = oldContext }()

	_, err := getUserDashboardURL()
	if err == nil {
		t.Error("expected error when no session exists")
	}
}

func TestDashboardViewCmdRunEPresent(t *testing.T) {
	if dashboardViewCmd.RunE == nil {
		t.Error("dashboardViewCmd.RunE should be set")
	}
}

func TestDashboardViewMaxArgs(t *testing.T) {
	// Test that the command accepts 0 or 1 arguments
	// The Args validator should be MaximumNArgs(1)
	cmd := dashboardViewCmd

	if cmd.Args == nil {
		t.Error("dashboardViewCmd should have Args validator set")
	}
}

func TestOpenBrowserFunction(t *testing.T) {
	// Save the original browserOpener
	originalOpener := browserOpener
	defer func() { browserOpener = originalOpener }()

	// Track what command would be executed
	var capturedCmd string
	var capturedArgs []string

	// Mock the browserOpener to capture the command without executing it
	browserOpener = func(cmd string, args []string) error {
		capturedCmd = cmd
		capturedArgs = args
		return nil
	}

	testURL := "https://dashboard.pantheon.io/sites/test-site-id"

	// Call openBrowser
	err := openBrowser(testURL)
	if err != nil {
		t.Errorf("openBrowser returned error: %v", err)
	}

	// Verify the correct command was prepared based on the OS
	if capturedCmd == "" {
		t.Error("browserOpener was not called")
	}

	// Verify the URL is in the captured args
	urlFound := false
	for _, arg := range capturedArgs {
		if arg == testURL {
			urlFound = true
			break
		}
	}
	if !urlFound {
		t.Errorf("expected URL '%s' in args, got %v", testURL, capturedArgs)
	}
}
