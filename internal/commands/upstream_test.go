package commands

import (
	"testing"
)

func TestUpstreamCmdStructure(t *testing.T) {
	if upstreamCmd.Use != "upstream" {
		t.Errorf("expected upstreamCmd.Use to be 'upstream', got '%s'", upstreamCmd.Use)
	}

	if upstreamCmd.Short == "" {
		t.Error("upstreamCmd.Short should not be empty")
	}
}

func TestUpstreamInfoCmdStructure(t *testing.T) {
	if upstreamInfoCmd.Use != "info <upstream>" {
		t.Errorf("expected upstreamInfoCmd.Use to be 'info <upstream>', got '%s'", upstreamInfoCmd.Use)
	}

	if upstreamInfoCmd.Short == "" {
		t.Error("upstreamInfoCmd.Short should not be empty")
	}

	// Verify upstreamInfoCmd is a subcommand of upstreamCmd
	found := false
	for _, cmd := range upstreamCmd.Commands() {
		if cmd.Name() == "info" {
			found = true
			break
		}
	}
	if !found {
		t.Error("upstreamInfoCmd should be a subcommand of upstreamCmd")
	}
}

func TestUpstreamListCmdStructure(t *testing.T) {
	if upstreamListCmd.Use != "list" {
		t.Errorf("expected upstreamListCmd.Use to be 'list', got '%s'", upstreamListCmd.Use)
	}

	if upstreamListCmd.Short == "" {
		t.Error("upstreamListCmd.Short should not be empty")
	}

	// Verify upstreamListCmd is a subcommand of upstreamCmd
	found := false
	for _, cmd := range upstreamCmd.Commands() {
		if cmd.Name() == "list" {
			found = true
			break
		}
	}
	if !found {
		t.Error("upstreamListCmd should be a subcommand of upstreamCmd")
	}

	// Verify aliases
	foundAlias := false
	for _, alias := range upstreamListCmd.Aliases {
		if alias == "ls" {
			foundAlias = true
			break
		}
	}
	if !foundAlias {
		t.Error("upstreamListCmd should have 'ls' as an alias")
	}
}

func TestUpstreamSubcommands(t *testing.T) {
	expectedSubcommands := []string{"info", "list"}
	subcommands := upstreamCmd.Commands()

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range subcommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand '%s' not found in upstreamCmd", expected)
		}
	}
}
