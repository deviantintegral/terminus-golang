package commands

import (
	"fmt"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/spf13/cobra"
)

var lockInfoCmd = &cobra.Command{
	Use:   "lock:info <site>.<env>",
	Short: "Show lock status",
	Long:  "Display HTTP basic authentication lock status for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runLockInfo,
}

var lockEnableCmd = &cobra.Command{
	Use:   "lock:enable <site>.<env>",
	Short: "Enable environment lock",
	Long:  "Enable HTTP basic authentication for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runLockEnable,
}

var lockDisableCmd = &cobra.Command{
	Use:   "lock:disable <site>.<env>",
	Short: "Disable environment lock",
	Long:  "Disable HTTP basic authentication for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runLockDisable,
}

var (
	lockUsername string
	lockPassword string
)

func init() {
	// Add lock commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(lockInfoCmd)
	rootCmd.AddCommand(lockEnableCmd)
	rootCmd.AddCommand(lockDisableCmd)

	// Enable flags
	lockEnableCmd.Flags().StringVar(&lockUsername, "username", "", "Username for HTTP basic auth")
	lockEnableCmd.Flags().StringVar(&lockPassword, "password", "", "Password for HTTP basic auth")
	_ = lockEnableCmd.MarkFlagRequired("username")
	_ = lockEnableCmd.MarkFlagRequired("password")
}

func runLockInfo(_ *cobra.Command, args []string) error {
	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	lock, err := envsService.GetLock(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to get lock info: %w", err)
	}

	return printOutput(lock)
}

func runLockEnable(_ *cobra.Command, args []string) error {
	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	printMessage("Enabling lock for %s.%s...", siteID, envID)

	if err := envsService.SetLock(getContext(), siteID, envID, lockUsername, lockPassword); err != nil {
		return fmt.Errorf("failed to enable lock: %w", err)
	}

	printMessage("Lock enabled successfully!")

	return nil
}

func runLockDisable(_ *cobra.Command, args []string) error {
	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	if !confirm(fmt.Sprintf("Are you sure you want to disable lock for %s.%s?", siteID, envID)) {
		printMessage("Canceled")
		return nil
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	printMessage("Disabling lock for %s.%s...", siteID, envID)

	if err := envsService.RemoveLock(getContext(), siteID, envID); err != nil {
		return fmt.Errorf("failed to disable lock: %w", err)
	}

	printMessage("Lock disabled successfully!")

	return nil
}
