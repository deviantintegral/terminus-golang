package commands

import (
	"testing"
)

func TestPlanInfoCmdStructure(t *testing.T) {
	if planInfoCmd.Use != "plan:info <site>" {
		t.Errorf("expected planInfoCmd.Use to be 'plan:info <site>', got '%s'", planInfoCmd.Use)
	}

	if planInfoCmd.Short == "" {
		t.Error("planInfoCmd.Short should not be empty")
	}
}

func TestPlanCommands(t *testing.T) {
	expectedCommands := []string{"plan:info"}

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
