package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:     "domain",
	Aliases: []string{"domains"},
	Short:   "Domain management commands",
	Long:    "Manage environment domains",
}

var domainListCmd = &cobra.Command{
	Use:     "list <site>.<env>",
	Aliases: []string{"ls"},
	Short:   "List domains",
	Long:    "Display a list of domains for an environment",
	Args:    cobra.ExactArgs(1),
	RunE:    runDomainList,
}

var domainAddCmd = &cobra.Command{
	Use:   "add <site>.<env> <domain>",
	Short: "Add a domain",
	Long:  "Add a domain to an environment",
	Args:  cobra.ExactArgs(2),
	RunE:  runDomainAdd,
}

var domainRemoveCmd = &cobra.Command{
	Use:   "remove <site>.<env> <domain>",
	Short: "Remove a domain",
	Long:  "Remove a domain from an environment",
	Args:  cobra.ExactArgs(2),
	RunE:  runDomainRemove,
}

var domainDNSCmd = &cobra.Command{
	Use:   "dns <site>.<env> <domain>",
	Short: "Show DNS recommendations",
	Long:  "Display DNS recommendations for a domain",
	Args:  cobra.ExactArgs(2),
	RunE:  runDomainDNS,
}

func init() {
	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainAddCmd)
	domainCmd.AddCommand(domainRemoveCmd)
	domainCmd.AddCommand(domainDNSCmd)
}

func runDomainList(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	domainsService := api.NewDomainsService(cliContext.APIClient)

	domains, err := domainsService.List(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to list domains: %w", err)
	}

	return printOutput(domains)
}

func runDomainAdd(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	domain := args[1]
	domainsService := api.NewDomainsService(cliContext.APIClient)

	printMessage("Adding domain %s to %s.%s...", domain, siteID, envID)

	addedDomain, err := domainsService.Add(getContext(), siteID, envID, domain)
	if err != nil {
		return fmt.Errorf("failed to add domain: %w", err)
	}

	printMessage("Domain added successfully!")

	return printOutput(addedDomain)
}

func runDomainRemove(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	domain := args[1]

	if !confirm(fmt.Sprintf("Are you sure you want to remove domain '%s' from %s.%s?", domain, siteID, envID)) {
		printMessage("Cancelled")
		return nil
	}

	domainsService := api.NewDomainsService(cliContext.APIClient)

	printMessage("Removing domain %s from %s.%s...", domain, siteID, envID)

	if err := domainsService.Remove(getContext(), siteID, envID, domain); err != nil {
		return fmt.Errorf("failed to remove domain: %w", err)
	}

	printMessage("Domain removed successfully!")

	return nil
}

func runDomainDNS(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	domain := args[1]
	domainsService := api.NewDomainsService(cliContext.APIClient)

	records, err := domainsService.GetDNS(getContext(), siteID, envID, domain)
	if err != nil {
		return fmt.Errorf("failed to get DNS records: %w", err)
	}

	return printOutput(records)
}
