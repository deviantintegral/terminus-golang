package commands

import (
	"testing"
)

func TestConnectionCmdStructure(t *testing.T) {
	if connectionCmd.Use != "connection" {
		t.Errorf("expected connectionCmd.Use to be 'connection', got '%s'", connectionCmd.Use)
	}

	if connectionCmd.Short == "" {
		t.Error("connectionCmd.Short should not be empty")
	}
}

func TestConnectionInfoCmdStructure(t *testing.T) {
	if connectionInfoCmd.Use != "info <site>.<env>" {
		t.Errorf("expected connectionInfoCmd.Use to be 'info <site>.<env>', got '%s'", connectionInfoCmd.Use)
	}

	if connectionInfoCmd.Short == "" {
		t.Error("connectionInfoCmd.Short should not be empty")
	}

	// Verify connectionInfoCmd is a subcommand of connectionCmd
	found := false
	for _, cmd := range connectionCmd.Commands() {
		if cmd.Name() == "info" {
			found = true
			break
		}
	}
	if !found {
		t.Error("connectionInfoCmd should be a subcommand of connectionCmd")
	}
}

func TestConnectionSetCmdStructure(t *testing.T) {
	if connectionSetCmd.Use != "set <site>.<env> <mode>" {
		t.Errorf("expected connectionSetCmd.Use to be 'set <site>.<env> <mode>', got '%s'", connectionSetCmd.Use)
	}

	if connectionSetCmd.Short == "" {
		t.Error("connectionSetCmd.Short should not be empty")
	}

	// Verify connectionSetCmd is a subcommand of connectionCmd
	found := false
	for _, cmd := range connectionCmd.Commands() {
		if cmd.Name() == "set" {
			found = true
			break
		}
	}
	if !found {
		t.Error("connectionSetCmd should be a subcommand of connectionCmd")
	}
}

func TestConnectionSubcommands(t *testing.T) {
	expectedSubcommands := []string{"info", "set"}
	subcommands := connectionCmd.Commands()

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range subcommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand '%s' not found in connectionCmd", expected)
		}
	}
}
