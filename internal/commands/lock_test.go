package commands

import (
	"testing"
)

func TestLockInfoCmdStructure(t *testing.T) {
	if lockInfoCmd.Use != "lock:info <site>.<env>" {
		t.Errorf("expected lockInfoCmd.Use to be 'lock:info <site>.<env>', got '%s'", lockInfoCmd.Use)
	}

	if lockInfoCmd.Short == "" {
		t.Error("lockInfoCmd.Short should not be empty")
	}
}

func TestLockEnableCmdStructure(t *testing.T) {
	if lockEnableCmd.Use != "lock:enable <site>.<env>" {
		t.Errorf("expected lockEnableCmd.Use to be 'lock:enable <site>.<env>', got '%s'", lockEnableCmd.Use)
	}

	if lockEnableCmd.Short == "" {
		t.Error("lockEnableCmd.Short should not be empty")
	}
}

func TestLockDisableCmdStructure(t *testing.T) {
	if lockDisableCmd.Use != "lock:disable <site>.<env>" {
		t.Errorf("expected lockDisableCmd.Use to be 'lock:disable <site>.<env>', got '%s'", lockDisableCmd.Use)
	}

	if lockDisableCmd.Short == "" {
		t.Error("lockDisableCmd.Short should not be empty")
	}
}

func TestLockCommands(t *testing.T) {
	expectedCommands := []string{"lock:info", "lock:enable", "lock:disable"}

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
