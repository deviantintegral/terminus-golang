package commands

import (
	"testing"
)

func TestEnvMetricsCmdStructure(t *testing.T) {
	if envMetricsCmd.Use != "env:metrics <site>[.<env>]" {
		t.Errorf("expected envMetricsCmd.Use to be 'env:metrics <site>[.<env>]', got '%s'", envMetricsCmd.Use)
	}

	if envMetricsCmd.Short == "" {
		t.Error("envMetricsCmd.Short should not be empty")
	}

	if envMetricsCmd.Long == "" {
		t.Error("envMetricsCmd.Long should not be empty")
	}

	// Verify the command has RunE
	if envMetricsCmd.RunE == nil {
		t.Error("envMetricsCmd.RunE should not be nil")
	}
}

func TestEnvMetricsFlags(t *testing.T) {
	// Test period flag
	periodFlag := envMetricsCmd.Flags().Lookup("period")
	if periodFlag == nil {
		t.Error("envMetricsCmd should have a 'period' flag")
		return
	}
	if periodFlag.DefValue != "day" {
		t.Errorf("expected period flag default to be 'day', got '%s'", periodFlag.DefValue)
	}

	// Test datapoints flag
	datapointsFlag := envMetricsCmd.Flags().Lookup("datapoints")
	if datapointsFlag == nil {
		t.Error("envMetricsCmd should have a 'datapoints' flag")
		return
	}
	if datapointsFlag.DefValue != "auto" {
		t.Errorf("expected datapoints flag default to be 'auto', got '%s'", datapointsFlag.DefValue)
	}
}

func TestEnvMetricsCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "env:metrics <site>[.<env>]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("env:metrics command should be registered with rootCmd")
	}
}

func TestBuildMetricsDuration(t *testing.T) {
	tests := []struct {
		name        string
		period      string
		datapoints  string
		expected    string
		expectError bool
	}{
		{
			name:        "day period with auto datapoints",
			period:      "day",
			datapoints:  "auto",
			expected:    "28d",
			expectError: false,
		},
		{
			name:        "week period with auto datapoints",
			period:      "week",
			datapoints:  "auto",
			expected:    "12w",
			expectError: false,
		},
		{
			name:        "month period with auto datapoints",
			period:      "month",
			datapoints:  "auto",
			expected:    "12m",
			expectError: false,
		},
		{
			name:        "day period with custom datapoints",
			period:      "day",
			datapoints:  "7",
			expected:    "7d",
			expectError: false,
		},
		{
			name:        "week period with custom datapoints",
			period:      "week",
			datapoints:  "4",
			expected:    "4w",
			expectError: false,
		},
		{
			name:        "month period with custom datapoints",
			period:      "month",
			datapoints:  "6",
			expected:    "6m",
			expectError: false,
		},
		{
			name:        "invalid period",
			period:      "year",
			datapoints:  "auto",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid datapoints - negative",
			period:      "day",
			datapoints:  "-1",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid datapoints - zero",
			period:      "day",
			datapoints:  "0",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid datapoints - non-numeric",
			period:      "day",
			datapoints:  "abc",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty period",
			period:      "",
			datapoints:  "auto",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildMetricsDuration(tt.period, tt.datapoints)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected '%s', got '%s'", tt.expected, result)
				}
			}
		})
	}
}

func TestFindEnvSeparator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "site.env format",
			input:    "mysite.dev",
			expected: 6,
		},
		{
			name:     "site only",
			input:    "mysite",
			expected: -1,
		},
		{
			name:     "uuid.env format",
			input:    "12345678-1234-1234-1234-123456789abc.live",
			expected: 36,
		},
		{
			name:     "empty string",
			input:    "",
			expected: -1,
		},
		{
			name:     "starts with period",
			input:    ".dev",
			expected: 0,
		},
		{
			name:     "multiple periods",
			input:    "site.name.dev",
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findEnvSeparator(tt.input)
			if result != tt.expected {
				t.Errorf("findEnvSeparator(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEnvListCmdStructure(t *testing.T) {
	if envListCmd.Use != "env:list <site>" {
		t.Errorf("expected envListCmd.Use to be 'env:list <site>', got '%s'", envListCmd.Use)
	}

	// Check alias
	foundAlias := false
	for _, alias := range envListCmd.Aliases {
		if alias == "environment" {
			foundAlias = true
			break
		}
	}
	if !foundAlias {
		t.Error("envListCmd should have 'environment' as an alias")
	}
}

func TestEnvInfoCmdStructure(t *testing.T) {
	if envInfoCmd.Use != "env:info <site>.<env>" {
		t.Errorf("expected envInfoCmd.Use to be 'env:info <site>.<env>', got '%s'", envInfoCmd.Use)
	}

	if envInfoCmd.Short == "" {
		t.Error("envInfoCmd.Short should not be empty")
	}
}

func TestEnvCommands(t *testing.T) {
	expectedCommands := []string{
		"env:list",
		"env:info",
		"env:clear-cache",
		"env:deploy",
		"env:clone-content",
		"env:commit",
		"env:wipe",
		"env:connection:set",
		"env:metrics",
	}

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
