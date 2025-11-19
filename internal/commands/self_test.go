package commands

import (
	"os"
	"testing"
)

func TestSelfInfoCmdStructure(t *testing.T) {
	if selfInfoCmd.Use != "self:info" {
		t.Errorf("expected selfInfoCmd.Use to be 'self:info', got '%s'", selfInfoCmd.Use)
	}

	if selfInfoCmd.Short == "" {
		t.Error("selfInfoCmd.Short should not be empty")
	}
}

func TestSelfCommands(t *testing.T) {
	expectedCommands := []string{"self:info"}

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
