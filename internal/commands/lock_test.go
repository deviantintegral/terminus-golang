package commands

import (
	"testing"
)

func TestLockCmdStructure(t *testing.T) {
	if lockCmd.Use != "lock" {
		t.Errorf("expected lockCmd.Use to be 'lock', got '%s'", lockCmd.Use)
	}

	if lockCmd.Short == "" {
		t.Error("lockCmd.Short should not be empty")
	}
}

func TestLockInfoCmdStructure(t *testing.T) {
	if lockInfoCmd.Use != "info <site>.<env>" {
		t.Errorf("expected lockInfoCmd.Use to be 'info <site>.<env>', got '%s'", lockInfoCmd.Use)
	}

	if lockInfoCmd.Short == "" {
		t.Error("lockInfoCmd.Short should not be empty")
	}

	// Verify lockInfoCmd is a subcommand of lockCmd
	found := false
	for _, cmd := range lockCmd.Commands() {
		if cmd.Name() == "info" {
			found = true
			break
		}
	}
	if !found {
		t.Error("lockInfoCmd should be a subcommand of lockCmd")
	}
}

func TestLockEnableCmdStructure(t *testing.T) {
	if lockEnableCmd.Use != "enable <site>.<env>" {
		t.Errorf("expected lockEnableCmd.Use to be 'enable <site>.<env>', got '%s'", lockEnableCmd.Use)
	}

	if lockEnableCmd.Short == "" {
		t.Error("lockEnableCmd.Short should not be empty")
	}

	// Verify lockEnableCmd is a subcommand of lockCmd
	found := false
	for _, cmd := range lockCmd.Commands() {
		if cmd.Name() == "enable" {
			found = true
			break
		}
	}
	if !found {
		t.Error("lockEnableCmd should be a subcommand of lockCmd")
	}
}

func TestLockDisableCmdStructure(t *testing.T) {
	if lockDisableCmd.Use != "disable <site>.<env>" {
		t.Errorf("expected lockDisableCmd.Use to be 'disable <site>.<env>', got '%s'", lockDisableCmd.Use)
	}

	if lockDisableCmd.Short == "" {
		t.Error("lockDisableCmd.Short should not be empty")
	}

	// Verify lockDisableCmd is a subcommand of lockCmd
	found := false
	for _, cmd := range lockCmd.Commands() {
		if cmd.Name() == "disable" {
			found = true
			break
		}
	}
	if !found {
		t.Error("lockDisableCmd should be a subcommand of lockCmd")
	}
}

func TestLockSubcommands(t *testing.T) {
	expectedSubcommands := []string{"info", "enable", "disable"}
	subcommands := lockCmd.Commands()

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range subcommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand '%s' not found in lockCmd", expected)
		}
	}
}

func TestLockEnableFlags(t *testing.T) {
	// Test that the username and password flags exist
	usernameFlag := lockEnableCmd.Flags().Lookup("username")
	if usernameFlag == nil {
		t.Error("lockEnableCmd should have a 'username' flag")
	}

	passwordFlag := lockEnableCmd.Flags().Lookup("password")
	if passwordFlag == nil {
		t.Error("lockEnableCmd should have a 'password' flag")
	}
}
