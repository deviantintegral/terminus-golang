package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:     "org",
	Aliases: []string{"organization"},
	Short:   "Organization management commands",
	Long:    "Manage organizations",
}

var orgListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List organizations",
	Long:    "Display a list of organizations",
	RunE:    runOrgList,
}

var orgInfoCmd = &cobra.Command{
	Use:   "info <org>",
	Short: "Show organization information",
	Long:  "Display detailed information about an organization",
	Args:  cobra.ExactArgs(1),
	RunE:  runOrgInfo,
}

var orgPeopleCmd = &cobra.Command{
	Use:   "people",
	Short: "Organization people management",
	Long:  "Manage organization members",
}

var orgPeopleListCmd = &cobra.Command{
	Use:     "list <org>",
	Aliases: []string{"ls"},
	Short:   "List organization members",
	Long:    "Display a list of members in an organization",
	Args:    cobra.ExactArgs(1),
	RunE:    runOrgPeopleList,
}

var orgSiteCmd = &cobra.Command{
	Use:   "site",
	Short: "Organization site management",
	Long:  "Manage organization sites",
}

var orgSiteListCmd = &cobra.Command{
	Use:     "list <org>",
	Aliases: []string{"ls"},
	Short:   "List organization sites",
	Long:    "Display a list of sites in an organization",
	Args:    cobra.ExactArgs(1),
	RunE:    runOrgSiteList,
}

var orgUpstreamsCmd = &cobra.Command{
	Use:   "upstreams",
	Short: "Organization upstream management",
	Long:  "Manage organization upstreams",
}

var orgUpstreamsListCmd = &cobra.Command{
	Use:     "list <org>",
	Aliases: []string{"ls"},
	Short:   "List organization upstreams",
	Long:    "Display a list of upstreams for an organization",
	Args:    cobra.ExactArgs(1),
	RunE:    runOrgUpstreamsList,
}

func init() {
	orgCmd.AddCommand(orgListCmd)
	orgCmd.AddCommand(orgInfoCmd)
	orgCmd.AddCommand(orgPeopleCmd)
	orgCmd.AddCommand(orgSiteCmd)
	orgCmd.AddCommand(orgUpstreamsCmd)

	orgPeopleCmd.AddCommand(orgPeopleListCmd)
	orgSiteCmd.AddCommand(orgSiteListCmd)
	orgUpstreamsCmd.AddCommand(orgUpstreamsListCmd)
}

func runOrgList(_ *cobra.Command, _ []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	orgsService := api.NewOrganizationsService(cliContext.APIClient)

	orgs, err := orgsService.List(getContext())
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
