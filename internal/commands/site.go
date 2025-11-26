package commands

import (
	"fmt"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/deviantintegral/terminus-golang/pkg/api/models"
	"github.com/spf13/cobra"
)

var siteListCmd = &cobra.Command{
	Use:     "site:list",
	Aliases: []string{"sites"},
	Short:   "List all sites",
	Long:    "Display a list of all sites accessible to the authenticated user",
	RunE:    runSiteList,
}

var siteInfoCmd = &cobra.Command{
	Use:   "site:info <site>",
	Short: "Show site information",
	Long:  "Display detailed information about a specific site",
	Args:  cobra.ExactArgs(1),
	RunE:  runSiteInfo,
}

var siteCreateCmd = &cobra.Command{
	Use:   "site:create <site_name> <label> <upstream_id>",
	Short: "Create a new site",
	Long:  "Creates a new site named <site_name>, human-readably labeled <label>, using code from <upstream_id>.",
	Args:  cobra.ExactArgs(3),
	RunE:  runSiteCreate,
}

var siteDeleteCmd = &cobra.Command{
	Use:   "site:delete <site>",
	Short: "Delete a site",
	Long:  "Delete a site from Pantheon",
	Args:  cobra.ExactArgs(1),
	RunE:  runSiteDelete,
}

var siteTeamListCmd = &cobra.Command{
	Use:   "site:team:list <site>",
	Short: "List site team members",
	Long:  "Display team members for a site",
	Args:  cobra.ExactArgs(1),
	RunE:  runSiteTeamList,
}

var siteOrgListCmd = &cobra.Command{
	Use:   "site:org:list <site>",
	Short: "List organizations for a site",
	Long:  "Display a list of supporting organizations for a site",
	Args:  cobra.ExactArgs(1),
	RunE:  runSiteOrgList,
}

var (
	siteOrgFlag    string
	siteRegionFlag string
)

func init() {
	// Add site commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(siteListCmd)
	rootCmd.AddCommand(siteInfoCmd)
	rootCmd.AddCommand(siteCreateCmd)
	rootCmd.AddCommand(siteDeleteCmd)
	rootCmd.AddCommand(siteTeamListCmd)
	rootCmd.AddCommand(siteOrgListCmd)

	// Flags
	siteListCmd.Flags().StringVar(&siteOrgFlag, "org", "", "Filter by organization")

	siteCreateCmd.Flags().StringVar(&siteOrgFlag, "org", "", "Organization ID")
	siteCreateCmd.Flags().StringVar(&siteRegionFlag, "region", "", "Preferred region")
}

func runSiteOrgList(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]
	sitesService := api.NewSitesService(cliContext.APIClient)

	orgs, err := sitesService.ListOrganizations(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to list site organizations: %w", err)
	}

	if len(orgs) == 0 {
		printMessage("This site has no supporting organizations")
		return nil
	}

	return printOutput(orgs)
}

func runSiteList(_ *cobra.Command, _ []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	// Load session to get user ID
	sess, err := cliContext.SessionStore.LoadSession()
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}
	if sess == nil || sess.UserID == "" {
		return fmt.Errorf("no user ID in session")
	}

	sitesService := api.NewSitesService(cliContext.APIClient)

	var sites []*models.Site

	if siteOrgFlag != "" {
		// If org flag is specified, only list sites for that specific organization
		sites, err = sitesService.ListByOrganization(getContext(), siteOrgFlag)
		if err != nil {
			return fmt.Errorf("failed to list sites: %w", err)
		}
	} else {
		// Otherwise, list all sites from user memberships and organization memberships
		sites, err = getAllUserSites(sess.UserID)
		if err != nil {
			return fmt.Errorf("failed to list sites: %w", err)
		}
	}

	// Convert to SiteListItem to exclude upstream field from output
	listItems := make([]*models.SiteListItem, len(sites))
	for i, site := range sites {
		listItems[i] = site.ToListItem()
	}

	return printOutput(listItems)
}

// getAllUserSites fetches all sites accessible to the user, including:
// 1. Sites from direct user memberships
// 2. Sites from all organizations the user is a member of
func getAllUserSites(userID string) ([]*models.Site, error) {
	sitesService := api.NewSitesService(cliContext.APIClient)
	orgsService := api.NewOrganizationsService(cliContext.APIClient)

	// Track unique sites by ID to avoid duplicates
	siteMap := make(map[string]*models.Site)

	// 1. Get sites from direct user memberships
	userSites, err := sitesService.List(getContext(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user sites: %w", err)
	}
	for _, site := range userSites {
		siteMap[site.ID] = site
	}

	// 2. Get user's organization memberships
	orgs, err := orgsService.List(getContext(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user organizations: %w", err)
	}

	// 3. For each organization, get all sites
	for _, org := range orgs {
		orgSites, err := sitesService.ListByOrganization(getContext(), org.ID)
		if err != nil {
			// Continue on error to get sites from other orgs
			// Log the error but don't fail completely
			orgName := org.ID
			if org.Label != "" {
				orgName = org.Label
			}
			printMessage("Warning: failed to list sites for organization %s: %v", orgName, err)
			continue
		}

		// Add org sites to the map (deduplicating by ID)
		for _, site := range orgSites {
			siteMap[site.ID] = site
		}
	}

	// Convert map to slice
	sites := make([]*models.Site, 0, len(siteMap))
	for _, site := range siteMap {
		sites = append(sites, site)
	}

	return sites, nil
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

	// Load session to get user ID
	sess, err := cliContext.SessionStore.LoadSession()
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}
	if sess == nil || sess.UserID == "" {
		return fmt.Errorf("no user ID in session")
	}

	siteName := args[0]
	label := args[1]
	upstreamID := args[2]
	sitesService := api.NewSitesService(cliContext.APIClient)

	req := &api.CreateSiteRequest{
		SiteName:     siteName,
		Label:        label,
		UpstreamID:   upstreamID,
		Organization: siteOrgFlag,
		Region:       siteRegionFlag,
	}

	printMessage("Creating site %s...", siteName)

	site, err := sitesService.Create(getContext(), sess.UserID, req)
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
		printMessage("Canceled")
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
