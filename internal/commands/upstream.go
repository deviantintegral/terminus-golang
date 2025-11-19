package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var upstreamInfoCmd = &cobra.Command{
	Use:   "upstream:info <upstream>",
	Short: "Show upstream information",
	Long:  "Display detailed information about an upstream",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpstreamInfo,
}

var upstreamListCmd = &cobra.Command{
	Use:   "upstream:list",
	Short: "List upstreams",
	Long:  "List available upstreams",
	Args:  cobra.NoArgs,
	RunE:  runUpstreamList,
}

func init() {
	// Add upstream commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(upstreamInfoCmd)
	rootCmd.AddCommand(upstreamListCmd)
}

func runUpstreamInfo(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	upstreamID := args[0]

	upstreamsService := api.NewUpstreamsService(cliContext.APIClient)

	upstream, err := upstreamsService.Get(getContext(), upstreamID)
	if err != nil {
		return fmt.Errorf("failed to get upstream info: %w", err)
	}

	return printOutput(upstream)
}

func runUpstreamList(_ *cobra.Command, _ []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	upstreamsService := api.NewUpstreamsService(cliContext.APIClient)

	upstreams, err := upstreamsService.List(getContext())
	if err != nil {
		return fmt.Errorf("failed to list upstreams: %w", err)
	}

	return printOutput(upstreams)
}
