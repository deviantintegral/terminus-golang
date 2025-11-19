package commands

import (
	"testing"
	"time"
)

func TestBackupListCmdStructure(t *testing.T) {
	// Test that backupListCmd has the expected properties
	if backupListCmd.Use != "backup:list <site>.<env>" {
		t.Errorf("expected backupListCmd.Use to be 'backup:list <site>.<env>', got '%s'", backupListCmd.Use)
	}

	if backupListCmd.Short == "" {
		t.Error("backupListCmd.Short should not be empty")
	}

	// Verify aliases
	foundAlias := false
	for _, alias := range backupListCmd.Aliases {
		if alias == "backups" {
			foundAlias = true
			break
		}
	}
	if !foundAlias {
		t.Error("backupListCmd should have 'backups' as an alias")
	}
}

func TestBackupInfoCmdStructure(t *testing.T) {
	if backupInfoCmd.Use != "backup:info <site>.<env>" {
		t.Errorf("expected backupInfoCmd.Use to be 'backup:info <site>.<env>', got '%s'", backupInfoCmd.Use)
	}

	if backupInfoCmd.Short == "" {
		t.Error("backupInfoCmd.Short should not be empty")
	}
}

func TestBackupCommands(t *testing.T) {
	expectedCommands := []string{"backup:list", "backup:create", "backup:get", "backup:restore", "backup:info", "backup:automatic:info", "backup:automatic:enable", "backup:automatic:disable"}

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
