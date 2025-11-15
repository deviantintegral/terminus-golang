package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:     "site",
	Aliases: []string{"sites"},
	Short:   "Site management commands",
	Long:    "Manage Pantheon sites",
}

var siteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all sites",
	Long:    "Display a list of all sites accessible to the authenticated user",
	RunE:    runSiteList,
}

var siteInfoCmd = &cobra.Command{
	Use:   "info <site>",
	Short: "Show site information",
	Long:  "Display detailed information about a specific site",
	Args:  cobra.ExactArgs(1),
	RunE:  runSiteInfo,
}

var siteCreateCmd = &cobra.Command{
	Use:   "create <site-name>",
	Short: "Create a new site",
	Long:  "Create a new site on Pantheon",
	Args:  cobra.ExactArgs(1),
	RunE:  runSiteCreate,
}

var siteDeleteCmd = &cobra.Command{
	Use:   "delete <site>",
	Short: "Delete a site",
	Long:  "Delete a site from Pantheon",
	Args:  cobra.ExactArgs(1),
	RunE:  runSiteDelete,
}

var siteTeamListCmd = &cobra.Command{
	Use:   "list <site>",
	Short: "List site team members",
	Long:  "Display team members for a site",
	Args:  cobra.ExactArgs(1),
	RunE:  runSiteTeamList,
}

var siteTeamCmd = &cobra.Command{
	Use:   "team",
	Short: "Site team management",
	Long:  "Manage site team members",
}

var (
	siteOrgFlag      string
	siteLabelFlag    string
	siteUpstreamFlag string
	siteRegionFlag   string
)

func init() {
	// Site commands
	siteCmd.AddCommand(siteListCmd)
	siteCmd.AddCommand(siteInfoCmd)
	siteCmd.AddCommand(siteCreateCmd)
	siteCmd.AddCommand(siteDeleteCmd)
	siteCmd.AddCommand(siteTeamCmd)

	// Team commands
	siteTeamCmd.AddCommand(siteTeamListCmd)

	// Flags
	siteListCmd.Flags().StringVar(&siteOrgFlag, "org", "", "Filter by organization")

	siteCreateCmd.Flags().StringVar(&siteLabelFlag, "label", "", "Site label")
	siteCreateCmd.Flags().StringVar(&siteUpstreamFlag, "upstream", "", "Upstream ID")
	siteCreateCmd.Flags().StringVar(&siteOrgFlag, "org", "", "Organization ID")
	siteCreateCmd.Flags().StringVar(&siteRegionFlag, "region", "", "Preferred region")
	_ = siteCreateCmd.MarkFlagRequired("upstream")
}

func runSiteList(_ *cobra.Command, _ []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	sitesService := api.NewSitesService(cliContext.APIClient)

	var sites interface{}
	var err error

	if siteOrgFlag != "" {
		sites, err = sitesService.ListByOrganization(getContext(), siteOrgFlag)
	} else {
		sites, err = sitesService.List(getContext())
	}

	if err != nil {
		return fmt.Errorf("failed to list sites: %w", err)
	}

	return printOutput(sites)
}

func runSiteInfo(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]
	sitesService := api.NewSitesService(cliContext.APIClient)

	site, err := sitesService.Get(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to get site info: %w", err)
	}

	return printOutput(site)
}

func runSiteCreate(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteName := args[0]
	sitesService := api.NewSitesService(cliContext.APIClient)

	req := &api.CreateSiteRequest{
		SiteName:     siteName,
		Label:        siteLabelFlag,
		UpstreamID:   siteUpstreamFlag,
		Organization: siteOrgFlag,
		Region:       siteRegionFlag,
	}

	if req.Label == "" {
		req.Label = siteName
	}

	printMessage("Creating site %s...", siteName)

	site, err := sitesService.Create(getContext(), req)
	if err != nil {
		return fmt.Errorf("failed to create site: %w", err)
	}

	printMessage("Site created successfully!")

	return printOutput(site)
}

func runSiteDelete(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]

	if !confirm(fmt.Sprintf("Are you sure you want to delete site '%s'? This cannot be undone.", siteID)) {
		printMessage("Cancelled")
		return nil
	}

	sitesService := api.NewSitesService(cliContext.APIClient)

	printMessage("Deleting site %s...", siteID)

	if err := sitesService.Delete(getContext(), siteID); err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}

	printMessage("Site deleted successfully")

	return nil
}

func runSiteTeamList(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]
	sitesService := api.NewSitesService(cliContext.APIClient)

	team, err := sitesService.GetTeam(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to get team members: %w", err)
	}

	return printOutput(team)
}
