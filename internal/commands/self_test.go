package commands

import (
	"os"
	"testing"
)

func TestSelfCmdStructure(t *testing.T) {
	if selfCmd.Use != "self" {
		t.Errorf("expected selfCmd.Use to be 'self', got '%s'", selfCmd.Use)
	}

	if selfCmd.Short == "" {
		t.Error("selfCmd.Short should not be empty")
	}
}

func TestSelfInfoCmdStructure(t *testing.T) {
	if selfInfoCmd.Use != "info" {
		t.Errorf("expected selfInfoCmd.Use to be 'info', got '%s'", selfInfoCmd.Use)
	}

	if selfInfoCmd.Short == "" {
		t.Error("selfInfoCmd.Short should not be empty")
	}

	// Verify selfInfoCmd is a subcommand of selfCmd
	found := false
	for _, cmd := range selfCmd.Commands() {
		if cmd.Name() == "info" {
			found = true
			break
		}
	}
	if !found {
		t.Error("selfInfoCmd should be a subcommand of selfCmd")
	}
}

func TestSelfSubcommands(t *testing.T) {
	expectedSubcommands := []string{"info"}
	subcommands := selfCmd.Commands()

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range subcommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand '%s' not found in selfCmd", expected)
		}
	}
}

func TestRunSelfInfo(t *testing.T) {
	// Initialize CLI context for printOutput
	cliContext = &CLIContext{}

	// Set quietFlag to avoid nil pointer in printOutput
	oldQuietFlag := quietFlag
	quietFlag = true
	defer func() { quietFlag = oldQuietFlag }()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSelfInfo(nil, nil)

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout
	_, _ = r.Read(make([]byte, 1024))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVersionVariable(t *testing.T) {
	// Test that Version variable is set
	if Version == "" {
		t.Error("Version should not be empty")
	}
}
