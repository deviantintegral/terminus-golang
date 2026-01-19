package commands

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/deviantintegral/terminus-golang/pkg/version"
	"github.com/spf13/cobra"
)

var selfInfoCmd = &cobra.Command{
	Use:   "self:info",
	Short: "Show Terminus information",
	Long:  "Display information about Terminus including version, paths, and runtime details",
	Args:  cobra.NoArgs,
	RunE:  runSelfInfo,
}

func init() {
	// Add self commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(selfInfoCmd)
}

func runSelfInfo(_ *cobra.Command, _ []string) error {
	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		execPath = "unknown"
	}

	// Get config directory
	homeDir, err := os.UserHomeDir()
	configPath := "unknown"
	if err == nil {
		configPath = filepath.Join(homeDir, ".terminus")
	}

	// Create info struct
	type SelfInfo struct {
		TerminusVersion string `json:"terminus_version"`
		TerminusPath    string `json:"terminus_path"`
		ConfigPath      string `json:"config_path"`
		GoVersion       string `json:"go_version"`
		OS              string `json:"os"`
		Architecture    string `json:"architecture"`
	}

	info := SelfInfo{
		TerminusVersion: version.String(),
		TerminusPath:    execPath,
		ConfigPath:      configPath,
		GoVersion:       runtime.Version(),
		OS:              runtime.GOOS,
		Architecture:    runtime.GOARCH,
	}

	return printOutput(info)
}
