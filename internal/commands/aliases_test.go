package commands

import (
	"strings"
	"testing"

	"github.com/deviantintegral/terminus-golang/pkg/api/models"
)

func TestAliasesCmdStructure(t *testing.T) {
	if aliasesCmd.Use != "aliases" {
		t.Errorf("expected aliasesCmd.Use to be 'aliases', got '%s'", aliasesCmd.Use)
	}

	if aliasesCmd.Short == "" {
		t.Error("aliasesCmd.Short should not be empty")
	}

	if aliasesCmd.Long == "" {
		t.Error("aliasesCmd.Long should not be empty")
	}

	// Check that it has the drush:aliases alias
	hasAlias := false
	for _, alias := range aliasesCmd.Aliases {
		if alias == "drush:aliases" {
			hasAlias = true
			break
		}
	}
	if !hasAlias {
		t.Error("aliasesCmd should have 'drush:aliases' alias")
	}

	// Check flags
	printFlag := aliasesCmd.Flags().Lookup("print")
	if printFlag == nil {
		t.Error("aliasesCmd should have a 'print' flag")
	}

	locationFlag := aliasesCmd.Flags().Lookup("location")
	if locationFlag == nil {
		t.Error("aliasesCmd should have a 'location' flag")
	}

	typeFlag := aliasesCmd.Flags().Lookup("type")
	if typeFlag == nil {
		t.Error("aliasesCmd should have a 'type' flag")
	}

	baseFlag := aliasesCmd.Flags().Lookup("base")
	if baseFlag == nil {
		t.Error("aliasesCmd should have a 'base' flag")
	}
}

func TestAliasesCommand(t *testing.T) {
	// Check that the aliases command is registered with rootCmd
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "aliases" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'aliases' command not found in rootCmd")
	}
}

func TestGenerateDrush8Aliases(t *testing.T) {
	sites := []*models.Site{
		{
			ID:   "11111111-1111-1111-1111-111111111111",
			Name: "test-site-1",
		},
		{
			ID:   "22222222-2222-2222-2222-222222222222",
			Name: "test-site-2",
		},
	}

	content := generateDrush8Aliases(sites)

	// Check header
	if !strings.Contains(content, "<?php") {
		t.Error("Drush 8 aliases should start with <?php")
	}

	if !strings.Contains(content, "Pantheon drush alias file") {
		t.Error("Drush 8 aliases should contain header comment")
	}

	// Check site 1 alias
	if !strings.Contains(content, "$aliases['test-site-1.*']") {
		t.Error("Drush 8 aliases should contain test-site-1 alias")
	}

	if !strings.Contains(content, "${env-name}-test-site-1.pantheonsite.io") {
		t.Error("Drush 8 aliases should contain URI for test-site-1")
	}

	if !strings.Contains(content, "appserver.${env-name}.11111111-1111-1111-1111-111111111111.drush.in") {
		t.Error("Drush 8 aliases should contain remote-host for test-site-1")
	}

	if !strings.Contains(content, "${env-name}.11111111-1111-1111-1111-111111111111") {
		t.Error("Drush 8 aliases should contain remote-user for test-site-1")
	}

	// Check site 2 alias
	if !strings.Contains(content, "$aliases['test-site-2.*']") {
		t.Error("Drush 8 aliases should contain test-site-2 alias")
	}

	// Check common elements
	if !strings.Contains(content, "'-p 2222 -o \"AddressFamily inet\"'") {
		t.Error("Drush 8 aliases should contain SSH options")
	}

	if !strings.Contains(content, "'%files' => 'files'") {
		t.Error("Drush 8 aliases should contain path-aliases")
	}
}

func TestGenerateDrush9Alias(t *testing.T) {
	site := &models.Site{
		ID:   "11111111-1111-1111-1111-111111111111",
		Name: "test-site",
	}

	content := generateDrush9Alias(site)

	// Check YAML structure
	if !strings.Contains(content, "'*':") {
		t.Error("Drush 9 alias should contain wildcard environment key")
	}

	if !strings.Contains(content, "host: appserver.${env-name}.11111111-1111-1111-1111-111111111111.drush.in") {
		t.Error("Drush 9 alias should contain host")
	}

	if !strings.Contains(content, "uri: ${env-name}-test-site.pantheonsite.io") {
		t.Error("Drush 9 alias should contain URI")
	}

	if !strings.Contains(content, "user: ${env-name}.11111111-1111-1111-1111-111111111111") {
		t.Error("Drush 9 alias should contain user")
	}

	if !strings.Contains(content, "paths:") {
		t.Error("Drush 9 alias should contain paths section")
	}

	if !strings.Contains(content, "files: files") {
		t.Error("Drush 9 alias should contain files path")
	}

	if !strings.Contains(content, "ssh:") {
		t.Error("Drush 9 alias should contain ssh section")
	}

	if !strings.Contains(content, "options: '-p 2222 -o \"AddressFamily inet\"'") {
		t.Error("Drush 9 alias should contain SSH options")
	}

	if !strings.Contains(content, "tty: false") {
		t.Error("Drush 9 alias should contain tty: false")
	}
}

func TestExpandHomePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "tilde path",
			input:    "~/.drush",
			contains: ".drush",
		},
		{
			name:     "absolute path",
			input:    "/absolute/path",
			contains: "/absolute/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandHomePath(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("expandHomePath(%s) = %s, should contain %s", tt.input, result, tt.contains)
			}
		})
	}
}

func TestShortenHomePath(t *testing.T) {
	// This test is environment-dependent, so we'll just verify it doesn't panic
	// and returns a non-empty string
	result := shortenHomePath("/some/absolute/path")
	if result == "" {
		t.Error("shortenHomePath should return a non-empty string")
	}
}
