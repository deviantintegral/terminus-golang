package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Redis management commands",
	Long:  "Manage Redis cache for Pantheon sites",
}

var redisEnableCmd = &cobra.Command{
	Use:   "enable <site>",
	Short: "Enable Redis for a site",
	Long:  "Enable Redis object cache for a Pantheon site",
	Args:  cobra.ExactArgs(1),
	RunE:  runRedisEnable,
}

var redisDisableCmd = &cobra.Command{
	Use:   "disable <site>",
	Short: "Disable Redis for a site",
	Long:  "Disable Redis object cache for a Pantheon site",
	Args:  cobra.ExactArgs(1),
	RunE:  runRedisDisable,
}

func init() {
	redisCmd.AddCommand(redisEnableCmd)
	redisCmd.AddCommand(redisDisableCmd)
}

func runRedisEnable(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]
	redisService := api.NewRedisService(cliContext.APIClient)

	printMessage("Enabling Redis for %s...", siteID)

	workflow, err := redisService.Enable(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to enable Redis: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Enabling Redis")
}

func runRedisDisable(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]
	redisService := api.NewRedisService(cliContext.APIClient)

	printMessage("Disabling Redis for %s...", siteID)

	workflow, err := redisService.Disable(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to disable Redis: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Disabling Redis")
}
