package commands

import (
	"fmt"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/spf13/cobra"
)

var tagListCmd = &cobra.Command{
	Use:   "tag:list <site> <org>",
	Short: "List tags for a site",
	Long:  "Display a list of tags for a site within an organization",
	Args:  cobra.ExactArgs(2),
	RunE:  runTagList,
}

func init() {
	rootCmd.AddCommand(tagListCmd)
}

func runTagList(_ *cobra.Command, args []string) error {
	siteID := args[0]
	orgID := args[1]

	sitesService := api.NewSitesService(cliContext.APIClient)

	tags, err := sitesService.GetTags(getContext(), siteID, orgID)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}

	if len(tags) == 0 {
		printMessage("No tags found for site %s in organization %s", siteID, orgID)
		return nil
	}

	return printOutput(tags)
}
