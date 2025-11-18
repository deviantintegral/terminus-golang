package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestArtCollection(t *testing.T) {
	// Verify all art pieces exist and have content
	expectedArt := []string{"fist", "hello", "rocket", "unicorn", "wordpress", "druplicon"}

	for _, name := range expectedArt {
		t.Run(name, func(t *testing.T) {
			art, exists := artCollection[name]
			if !exists {
				t.Errorf("art '%s' not found in artCollection", name)
				return
			}
			if art == "" {
				t.Errorf("art '%s' has empty content", name)
			}

			// Verify description exists
			desc, hasDesc := artDescriptions[name]
			if !hasDesc {
				t.Errorf("art '%s' has no description", name)
			}
			if desc == "" {
				t.Errorf("art '%s' has empty description", name)
			}
		})
	}
}

func TestArtCollectionCount(t *testing.T) {
	// Verify the expected number of art pieces
	expectedCount := 6
	actualCount := len(artCollection)

	if actualCount != expectedCount {
		t.Errorf("expected %d art pieces, got %d", expectedCount, actualCount)
	}

	// Verify descriptions match collection
	if len(artDescriptions) != len(artCollection) {
		t.Errorf("mismatch between artCollection (%d) and artDescriptions (%d)",
			len(artCollection), len(artDescriptions))
	}
}

func TestGetRandomArtName(t *testing.T) {
	// Run multiple times to ensure we get valid names
	for i := 0; i < 100; i++ {
		name := getRandomArtName()
		if _, exists := artCollection[name]; !exists {
			t.Errorf("getRandomArtName returned invalid name: %s", name)
		}
	}
}

func TestRunArt(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		contains    string
	}{
		{
			name:     "display fist art",
			args:     []string{"fist"},
			wantErr:  false,
			contains: "50 41 4E 54 48 45 4F 4E",
		},
		{
			name:     "display hello art",
			args:     []string{"hello"},
			wantErr:  false,
			contains: "Hello World!",
		},
		{
			name:     "display rocket art",
			args:     []string{"rocket"},
			wantErr:  false,
			contains: "/___\\",
		},
		{
			name:     "display unicorn art",
			args:     []string{"unicorn"},
			wantErr:  false,
			contains: ">\\/7",
		},
		{
			name:     "display wordpress art",
			args:     []string{"wordpress"},
			wantErr:  false,
			contains: "...............",
		},
		{
			name:     "display druplicon art",
			args:     []string{"druplicon"},
			wantErr:  false,
			contains: "MMMMMMMM",
		},
		{
			name:        "invalid art name",
			args:        []string{"nonexistent"},
			wantErr:     true,
			errContains: "art 'nonexistent' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := runArt(nil, tt.args)

			// Restore stdout and read captured output
			_ = w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain '%s', got: %s", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if tt.contains != "" && !strings.Contains(output, tt.contains) {
					t.Errorf("expected output to contain '%s', got: %s", tt.contains, output)
				}
			}
		})
	}
}

func TestRunArtRandom(t *testing.T) {
	// Test with no arguments (random selection)
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runArt(nil, []string{})

	_ = w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Errorf("unexpected error for random art: %v", err)
	}
	if output == "" {
		t.Error("expected some output for random art, got empty string")
	}
}

func TestRunArtList(t *testing.T) {
	// Initialize CLI context for printOutput
	cliContext = &CLIContext{}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set quietFlag to avoid nil pointer in printOutput
	oldQuietFlag := quietFlag
	quietFlag = true
	defer func() { quietFlag = oldQuietFlag }()

	err := runArtList(nil, nil)

	_ = w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestArtCmdStructure(t *testing.T) {
	// Test command structure
	if artCmd.Use != "art [name]" {
		t.Errorf("expected artCmd.Use to be 'art [name]', got '%s'", artCmd.Use)
	}

	if artCmd.Short == "" {
		t.Error("artCmd.Short should not be empty")
	}

	// Test that artListCmd is a subcommand
	found := false
	for _, cmd := range artCmd.Commands() {
		if cmd.Name() == "list" {
			found = true
			break
		}
	}
	if !found {
		t.Error("artListCmd should be a subcommand of artCmd")
	}
}

func TestArtListCmdStructure(t *testing.T) {
	if artListCmd.Use != "list" {
		t.Errorf("expected artListCmd.Use to be 'list', got '%s'", artListCmd.Use)
	}

	if artListCmd.Short == "" {
		t.Error("artListCmd.Short should not be empty")
	}
}

func TestArtInfoStruct(t *testing.T) {
	// Test that ArtInfo can be created and fields are accessible
	info := ArtInfo{
		Name:        "test",
		Description: "test description",
	}

	if info.Name != "test" {
		t.Errorf("expected Name to be 'test', got '%s'", info.Name)
	}
	if info.Description != "test description" {
		t.Errorf("expected Description to be 'test description', got '%s'", info.Description)
	}
}

func TestAllArtHasDescriptions(t *testing.T) {
	// Verify every art piece has a corresponding description
	for name := range artCollection {
		if _, exists := artDescriptions[name]; !exists {
			t.Errorf("art '%s' is missing a description", name)
		}
	}

	// Verify no orphaned descriptions
	for name := range artDescriptions {
		if _, exists := artCollection[name]; !exists {
			t.Errorf("description for '%s' has no corresponding art", name)
		}
	}
}
