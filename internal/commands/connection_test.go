package commands

import (
	"testing"
)

func TestConnectionInfoCmdStructure(t *testing.T) {
	if connectionInfoCmd.Use != "connection:info <site>.<env>" {
		t.Errorf("expected connectionInfoCmd.Use to be 'connection:info <site>.<env>', got '%s'", connectionInfoCmd.Use)
	}

	if connectionInfoCmd.Short == "" {
		t.Error("connectionInfoCmd.Short should not be empty")
	}
}

func TestConnectionSetCmdStructure(t *testing.T) {
	if connectionSetCmd.Use != "connection:set <site>.<env> <mode>" {
		t.Errorf("expected connectionSetCmd.Use to be 'connection:set <site>.<env> <mode>', got '%s'", connectionSetCmd.Use)
	}

	if connectionSetCmd.Short == "" {
		t.Error("connectionSetCmd.Short should not be empty")
	}
}

func TestConnectionCommands(t *testing.T) {
	expectedCommands := []string{"connection:info", "connection:set"}

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
