package commands

import (
	"fmt"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/spf13/cobra"
)

var branchListCmd = &cobra.Command{
	Use:   "branch:list <site>",
	Short: "List git branches for a site",
	Long:  "Display a list of git branches for a site",
	Args:  cobra.ExactArgs(1),
	RunE:  runBranchList,
}

func init() {
	rootCmd.AddCommand(branchListCmd)
}

func runBranchList(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]
	sitesService := api.NewSitesService(cliContext.APIClient)

	branches, err := sitesService.ListBranches(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	if len(branches) == 0 {
		printMessage("No branches found for site %s", siteID)
		return nil
	}

	return printOutput(branches)
}
