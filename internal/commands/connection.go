package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var connectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "Connection management commands",
	Long:  "Manage environment connection settings and display connection information",
}

var connectionInfoCmd = &cobra.Command{
	Use:   "info <site>.<env>",
	Short: "Show connection information",
	Long:  "Display connection information for an environment including SFTP, Git, MySQL, and Redis",
	Args:  cobra.ExactArgs(1),
	RunE:  runConnectionInfo,
}

var connectionSetCmd = &cobra.Command{
	Use:   "set <site>.<env> <mode>",
	Short: "Set connection mode",
	Long:  "Set the connection mode (git or sftp) for an environment",
	Args:  cobra.ExactArgs(2),
	RunE:  runConnectionSet,
}

func init() {
	connectionCmd.AddCommand(connectionInfoCmd)
	connectionCmd.AddCommand(connectionSetCmd)
}

func runConnectionInfo(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	info, err := envsService.GetConnectionInfo(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to get connection info: %w", err)
	}

	return printOutput(info)
}

func runConnectionSet(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	mode := args[1]
	if mode != "git" && mode != "sftp" {
		return fmt.Errorf("invalid connection mode: %s (must be 'git' or 'sftp')", mode)
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	printMessage("Setting connection mode to %s for %s.%s...", mode, siteID, envID)

	workflow, err := envsService.ChangeConnectionMode(getContext(), siteID, envID, mode)
	if err != nil {
		return fmt.Errorf("failed to set connection mode: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Setting connection mode")
}
