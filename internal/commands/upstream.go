package commands

import (
	"fmt"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/deviantintegral/terminus-golang/pkg/api/models"
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

var upstreamUpdatesListCmd = &cobra.Command{
	Use:   "upstream:updates:list <site>.<env>",
	Short: "List upstream updates",
	Long:  "Display a list of available upstream updates for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpstreamUpdatesList,
}

var (
	upstreamOrgFlag       string
	upstreamFrameworkFlag string
	upstreamAllFlag       bool
)

func init() {
	// Add upstream commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(upstreamInfoCmd)
	rootCmd.AddCommand(upstreamListCmd)
	rootCmd.AddCommand(upstreamUpdatesListCmd)

	upstreamListCmd.Flags().StringVar(&upstreamOrgFlag, "org", "", "Filter by organization")
	upstreamListCmd.Flags().StringVar(&upstreamFrameworkFlag, "framework", "", "Filter by framework")
	upstreamListCmd.Flags().BoolVar(&upstreamAllFlag, "all", false, "Show all upstreams")
}

func runUpstreamInfo(_ *cobra.Command, args []string) error {
	upstreamID := args[0]

	upstreamsService := api.NewUpstreamsService(cliContext.APIClient)

	upstream, err := upstreamsService.Get(getContext(), upstreamID)
	if err != nil {
		return fmt.Errorf("failed to get upstream info: %w", err)
	}

	return printOutput(upstream)
}

func runUpstreamList(_ *cobra.Command, _ []string) error {
	// Load session to get user ID
	sess, err := cliContext.SessionStore.LoadSession()
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}
	if sess == nil || sess.UserID == "" {
		return fmt.Errorf("no user ID in session")
	}

	upstreamsService := api.NewUpstreamsService(cliContext.APIClient)

	upstreams, err := upstreamsService.List(getContext(), sess.UserID)
	if err != nil {
		return fmt.Errorf("failed to list upstreams: %w", err)
	}

	upstreams = filterUpstreams(upstreams)

	if len(upstreams) == 0 {
		printMessage("No upstreams found")
		return nil
	}

	return printOutput(upstreams)
}

// filterUpstreams filters upstreams based on command flags
func filterUpstreams(upstreams []*models.Upstream) []*models.Upstream {
	// Filter by organization if specified
	if upstreamOrgFlag != "" {
		upstreams = filterUpstreamsByOrg(upstreams, upstreamOrgFlag)
	}

	// Filter by framework if specified
	if upstreamFrameworkFlag != "" {
		upstreams = filterUpstreamsByFramework(upstreams, upstreamFrameworkFlag)
	}

	// Filter to core and custom types unless --all is specified
	if !upstreamAllFlag {
		upstreams = filterUpstreamsByType(upstreams)
	}

	return upstreams
}

// filterUpstreamsByOrg filters upstreams by organization
func filterUpstreamsByOrg(upstreams []*models.Upstream, org string) []*models.Upstream {
	filtered := make([]*models.Upstream, 0)
	for _, u := range upstreams {
		if u.Organization == org {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

// filterUpstreamsByFramework filters upstreams by framework
func filterUpstreamsByFramework(upstreams []*models.Upstream, framework string) []*models.Upstream {
	filtered := make([]*models.Upstream, 0)
	for _, u := range upstreams {
		if u.Framework == framework {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

// filterUpstreamsByType filters upstreams to only core and custom types
func filterUpstreamsByType(upstreams []*models.Upstream) []*models.Upstream {
	filtered := make([]*models.Upstream, 0)
	for _, u := range upstreams {
		if u.Type == "core" || u.Type == "custom" {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func runUpstreamUpdatesList(_ *cobra.Command, args []string) error {
	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	upstreamsService := api.NewUpstreamsService(cliContext.APIClient)

	updates, err := upstreamsService.ListUpdates(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to list upstream updates: %w", err)
	}

	if len(updates) == 0 {
		printMessage("No upstream updates available")
		return nil
	}

	return printOutput(updates)
}
