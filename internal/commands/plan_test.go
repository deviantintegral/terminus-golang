package commands

import (
	"testing"
)

func TestPlanCmdStructure(t *testing.T) {
	if planCmd.Use != "plan" {
		t.Errorf("expected planCmd.Use to be 'plan', got '%s'", planCmd.Use)
	}

	if planCmd.Short == "" {
		t.Error("planCmd.Short should not be empty")
	}
}

func TestPlanInfoCmdStructure(t *testing.T) {
	if planInfoCmd.Use != "info <site>" {
		t.Errorf("expected planInfoCmd.Use to be 'info <site>', got '%s'", planInfoCmd.Use)
	}

	if planInfoCmd.Short == "" {
		t.Error("planInfoCmd.Short should not be empty")
	}

	// Verify planInfoCmd is a subcommand of planCmd
	found := false
	for _, cmd := range planCmd.Commands() {
		if cmd.Name() == "info" {
			found = true
			break
		}
	}
	if !found {
		t.Error("planInfoCmd should be a subcommand of planCmd")
	}
}

func TestPlanSubcommands(t *testing.T) {
	expectedSubcommands := []string{"info"}
	subcommands := planCmd.Commands()

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range subcommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand '%s' not found in planCmd", expected)
		}
	}
}
