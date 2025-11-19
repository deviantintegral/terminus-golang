package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var orgListCmd = &cobra.Command{
	Use:     "org:list",
	Aliases: []string{"organization"},
	Short:   "List organizations",
	Long:    "Display a list of organizations",
	RunE:    runOrgList,
}

var orgInfoCmd = &cobra.Command{
	Use:   "org:info <org>",
	Short: "Show organization information",
	Long:  "Display detailed information about an organization",
	Args:  cobra.ExactArgs(1),
	RunE:  runOrgInfo,
}

var orgPeopleListCmd = &cobra.Command{
	Use:   "org:people:list <org>",
	Short: "List organization members",
	Long:  "Display a list of members in an organization",
	Args:  cobra.ExactArgs(1),
	RunE:  runOrgPeopleList,
}

var orgSiteListCmd = &cobra.Command{
	Use:   "org:site:list <org>",
	Short: "List organization sites",
	Long:  "Display a list of sites in an organization",
	Args:  cobra.ExactArgs(1),
	RunE:  runOrgSiteList,
}

var orgUpstreamsListCmd = &cobra.Command{
	Use:   "org:upstreams:list <org>",
	Short: "List organization upstreams",
	Long:  "Display a list of upstreams for an organization",
	Args:  cobra.ExactArgs(1),
	RunE:  runOrgUpstreamsList,
}

func init() {
	// Add org commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(orgListCmd)
	rootCmd.AddCommand(orgInfoCmd)
	rootCmd.AddCommand(orgPeopleListCmd)
	rootCmd.AddCommand(orgSiteListCmd)
	rootCmd.AddCommand(orgUpstreamsListCmd)
}

func runOrgList(_ *cobra.Command, _ []string) error {
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

	orgsService := api.NewOrganizationsService(cliContext.APIClient)

	orgs, err := orgsService.List(getContext(), sess.UserID)
	if err != nil {
		return fmt.Errorf("failed to list organizations: %w", err)
	}

	return printOutput(orgs)
}

func runOrgInfo(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	orgID := args[0]
	orgsService := api.NewOrganizationsService(cliContext.APIClient)

	org, err := orgsService.Get(getContext(), orgID)
	if err != nil {
		return fmt.Errorf("failed to get organization info: %w", err)
	}

	return printOutput(org)
}

func runOrgPeopleList(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	orgID := args[0]
	orgsService := api.NewOrganizationsService(cliContext.APIClient)

	members, err := orgsService.ListMembers(getContext(), orgID)
	if err != nil {
		return fmt.Errorf("failed to list organization members: %w", err)
	}

	return printOutput(members)
}

func runOrgSiteList(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	orgID := args[0]
	sitesService := api.NewSitesService(cliContext.APIClient)

	sites, err := sitesService.ListByOrganization(getContext(), orgID)
	if err != nil {
		return fmt.Errorf("failed to list organization sites: %w", err)
	}

	return printOutput(sites)
}

func runOrgUpstreamsList(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	orgID := args[0]
	orgsService := api.NewOrganizationsService(cliContext.APIClient)

	upstreams, err := orgsService.ListUpstreams(getContext(), orgID)
	if err != nil {
		return fmt.Errorf("failed to list upstreams: %w", err)
	}

	return printOutput(upstreams)
}
