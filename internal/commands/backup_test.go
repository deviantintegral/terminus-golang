package commands

import (
	"testing"
	"time"
)

func TestBackupCmdStructure(t *testing.T) {
	// Test that backupCmd has the expected properties
	if backupCmd.Use != "backup" {
		t.Errorf("expected backupCmd.Use to be 'backup', got '%s'", backupCmd.Use)
	}

	if backupCmd.Short == "" {
		t.Error("backupCmd.Short should not be empty")
	}

	// Verify aliases
	foundAlias := false
	for _, alias := range backupCmd.Aliases {
		if alias == "backups" {
			foundAlias = true
			break
		}
	}
	if !foundAlias {
		t.Error("backupCmd should have 'backups' as an alias")
	}
}

func TestBackupInfoCmdStructure(t *testing.T) {
	if backupInfoCmd.Use != "info <site>.<env>" {
		t.Errorf("expected backupInfoCmd.Use to be 'info <site>.<env>', got '%s'", backupInfoCmd.Use)
	}

	if backupInfoCmd.Short == "" {
		t.Error("backupInfoCmd.Short should not be empty")
	}

	// Verify backupInfoCmd is a subcommand of backupCmd
	found := false
	for _, cmd := range backupCmd.Commands() {
		if cmd.Name() == "info" {
			found = true
			break
		}
	}
	if !found {
		t.Error("backupInfoCmd should be a subcommand of backupCmd")
	}
}

func TestBackupSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list", "create", "get", "restore", "info", "automatic"}
	subcommands := backupCmd.Commands()

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range subcommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand '%s' not found in backupCmd", expected)
		}
	}
}

func TestBackupAutomaticSubcommands(t *testing.T) {
	expectedSubcommands := []string{"info", "enable", "disable"}
	subcommands := backupAutomaticCmd.Commands()

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range subcommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand '%s' not found in backupAutomaticCmd", expected)
		}
	}
}

func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		expected  string
	}{
		{
			name:      "zero timestamp",
			timestamp: 0,
			expected:  "Never",
		},
		{
			name:      "specific timestamp",
			timestamp: 1609459200, // 2021-01-01 00:00:00 UTC
			expected:  time.Unix(1609459200, 0).Format("2006-01-02 15:04:05"),
		},
		{
			name:      "recent timestamp",
			timestamp: 1700000000,
			expected:  time.Unix(1700000000, 0).Format("2006-01-02 15:04:05"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimestamp(tt.timestamp)
			if result != tt.expected {
				t.Errorf("formatTimestamp(%d) = %s, want %s", tt.timestamp, result, tt.expected)
			}
		})
	}
}

func TestBackupInfoFlags(t *testing.T) {
	// Test that the element flag exists
	flag := backupInfoCmd.Flags().Lookup("element")
	if flag == nil {
		t.Error("backupInfoCmd should have an 'element' flag")
		return
	}

	if flag.DefValue != "files" {
		t.Errorf("expected element flag default to be 'files', got '%s'", flag.DefValue)
	}
}
