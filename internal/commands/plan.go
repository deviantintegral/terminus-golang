package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var planInfoCmd = &cobra.Command{
	Use:   "plan:info <site>",
	Short: "Show site plan information",
	Long:  "Display the current service plan for a site",
	Args:  cobra.ExactArgs(1),
	RunE:  runPlanInfo,
}

func init() {
	// Add plan commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(planInfoCmd)
}

func runPlanInfo(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]

	sitesService := api.NewSitesService(cliContext.APIClient)

	plan, err := sitesService.GetPlan(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to get site plan: %w", err)
	}

	return printOutput(plan)
}
