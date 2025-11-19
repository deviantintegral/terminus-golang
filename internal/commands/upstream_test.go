package commands

import (
	"testing"
)

func TestUpstreamInfoCmdStructure(t *testing.T) {
	if upstreamInfoCmd.Use != "upstream:info <upstream>" {
		t.Errorf("expected upstreamInfoCmd.Use to be 'upstream:info <upstream>', got '%s'", upstreamInfoCmd.Use)
	}

	if upstreamInfoCmd.Short == "" {
		t.Error("upstreamInfoCmd.Short should not be empty")
	}
}

func TestUpstreamListCmdStructure(t *testing.T) {
	if upstreamListCmd.Use != "upstream:list" {
		t.Errorf("expected upstreamListCmd.Use to be 'upstream:list', got '%s'", upstreamListCmd.Use)
	}

	if upstreamListCmd.Short == "" {
		t.Error("upstreamListCmd.Short should not be empty")
	}
}

func TestUpstreamCommands(t *testing.T) {
	expectedCommands := []string{"upstream:info", "upstream:list"}

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
