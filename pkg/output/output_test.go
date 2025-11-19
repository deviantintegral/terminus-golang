// Package output provides formatting utilities for CLI output.
package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintVerticalTable(t *testing.T) {
	// Create a test user-like struct
	type TestUser struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}

	user := TestUser{
		ID:        "dbd0f8fa-cd5f-4c55-9292-73c2e58edeb1",
		Email:     "andrew.berry@lullabot.com",
		FirstName: "Andrew",
		LastName:  "Berry",
	}

	var buf bytes.Buffer
	opts := &Options{
		Format: FormatTable,
		Writer: &buf,
	}

	err := Print(user, opts)
	if err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()
	t.Logf("Output:\n%s", output)

	// Verify output contains expected elements
	expectedElements := []string{
		"ID",
		"Email",
		"First Name",
		"Last Name",
		"dbd0f8fa-cd5f-4c55-9292-73c2e58edeb1",
		"andrew.berry@lullabot.com",
		"Andrew",
		"Berry",
		"---", // Border dashes
	}

	for _, elem := range expectedElements {
		if !strings.Contains(output, elem) {
			t.Errorf("Output missing expected element: %s", elem)
		}
	}

	// Verify vertical layout (should have multiple lines)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 6 {
		t.Errorf("Expected at least 6 lines (borders + 4 fields), got %d", len(lines))
	}

	// Verify top and bottom borders exist
	if !strings.Contains(lines[0], "---") {
		t.Error("Expected top border with dashes")
	}
	if !strings.Contains(lines[len(lines)-1], "---") {
		t.Error("Expected bottom border with dashes")
	}
}

func TestToHumanReadable(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"id", "ID"},
		{"firstname", "First Name"},
		{"lastname", "Last Name"},
		{"email", "Email"},
		{"profile", "Profile"},
		{"userName", "User Name"},
		{"user_name", "User Name"},
		{"createdAt", "Created At"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toHumanReadable(tt.input)
			if result != tt.expected {
				t.Errorf("toHumanReadable(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPrintHorizontalTable(t *testing.T) {
	// Test that multiple items still use horizontal layout
	type TestItem struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	items := []TestItem{
		{Name: "item1", Value: "value1"},
		{Name: "item2", Value: "value2"},
	}

	var buf bytes.Buffer
	opts := &Options{
		Format: FormatTable,
		Writer: &buf,
	}

	err := Print(items, opts)
	if err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()
	t.Logf("Output:\n%s", output)

	// Verify horizontal layout (headers should be on first line)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines (header + separator + rows), got %d", len(lines))
	}

	// First line should contain both field names
	firstLine := lines[0]
	if !strings.Contains(firstLine, "name") && !strings.Contains(firstLine, "Name") {
		t.Error("Expected first line to contain 'name' or 'Name' header")
	}
	if !strings.Contains(firstLine, "value") && !strings.Contains(firstLine, "Value") {
		t.Error("Expected first line to contain 'value' or 'Value' header")
	}
}
