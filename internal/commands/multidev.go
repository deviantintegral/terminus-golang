package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var multidevCreateCmd = &cobra.Command{
	Use:   "multidev:create <site>.<multidev>",
	Short: "Create a multidev environment",
	Long:  "Create a new multidev environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runMultidevCreate,
}

var multidevDeleteCmd = &cobra.Command{
	Use:   "multidev:delete <site>.<multidev>",
	Short: "Delete a multidev environment",
	Long:  "Delete a multidev environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runMultidevDelete,
}

var multidevMergeToDevCmd = &cobra.Command{
	Use:   "multidev:merge-to-dev <site>.<multidev>",
	Short: "Merge multidev to dev",
	Long:  "Merge a multidev environment into dev",
	Args:  cobra.ExactArgs(1),
	RunE:  runMultidevMergeToDev,
}

var multidevMergeFromDevCmd = &cobra.Command{
	Use:   "multidev:merge-from-dev <site>.<multidev>",
	Short: "Merge dev into multidev",
	Long:  "Merge dev into a multidev environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runMultidevMergeFromDev,
}

var multidevListCmd = &cobra.Command{
	Use:   "multidev:list <site>",
	Short: "List multidev environments",
	Long:  "Display a list of multidev environments for a site",
	Args:  cobra.ExactArgs(1),
	RunE:  runMultidevList,
}

var (
	multidevFromEnvFlag      string
	multidevDeleteBranchFlag bool
)

func init() {
	// Add multidev commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(multidevCreateCmd)
	rootCmd.AddCommand(multidevDeleteCmd)
	rootCmd.AddCommand(multidevMergeToDevCmd)
	rootCmd.AddCommand(multidevMergeFromDevCmd)
	rootCmd.AddCommand(multidevListCmd)

	multidevCreateCmd.Flags().StringVar(&multidevFromEnvFlag, "from-env", "dev", "Source environment")
	multidevDeleteCmd.Flags().BoolVar(&multidevDeleteBranchFlag, "delete-branch", false, "Delete git branch")
	multidevMergeToDevCmd.Flags().BoolVar(&envUpdateDBFlag, "updatedb", false, "Run database updates after merge")
	multidevMergeFromDevCmd.Flags().BoolVar(&envUpdateDBFlag, "updatedb", false, "Run database updates after merge")
}

func runMultidevList(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]
	multidevService := api.NewMultidevService(cliContext.APIClient)

	multidevs, err := multidevService.List(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to list multidev environments: %w", err)
	}

	if len(multidevs) == 0 {
		printMessage("You have no multidev environments")
		return nil
	}

	return printOutput(multidevs)
}

func runMultidevCreate(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envName, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	multidevService := api.NewMultidevService(cliContext.APIClient)

	printMessage("Creating multidev %s from %s...", envName, multidevFromEnvFlag)

	workflow, err := multidevService.Create(getContext(), siteID, envName, multidevFromEnvFlag)
	if err != nil {
		return fmt.Errorf("failed to create multidev: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Creating multidev")
}

func runMultidevDelete(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	if !confirm(fmt.Sprintf("Are you sure you want to delete multidev %s.%s?", siteID, envID)) {
		printMessage("Canceled")
		return nil
	}

	multidevService := api.NewMultidevService(cliContext.APIClient)

	printMessage("Deleting multidev %s.%s...", siteID, envID)

	workflow, err := multidevService.Delete(getContext(), siteID, envID, multidevDeleteBranchFlag)
	if err != nil {
		return fmt.Errorf("failed to delete multidev: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Deleting multidev")
}

func runMultidevMergeToDev(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	multidevService := api.NewMultidevService(cliContext.APIClient)

	printMessage("Merging %s.%s to dev...", siteID, envID)

	workflow, err := multidevService.MergeToDev(getContext(), siteID, envID, envUpdateDBFlag)
	if err != nil {
		return fmt.Errorf("failed to merge to dev: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Merging to dev")
}

func runMultidevMergeFromDev(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	multidevService := api.NewMultidevService(cliContext.APIClient)

	printMessage("Merging dev into %s.%s...", siteID, envID)

	workflow, err := multidevService.MergeFromDev(getContext(), siteID, envID, envUpdateDBFlag)
	if err != nil {
		return fmt.Errorf("failed to merge from dev: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Merging from dev")
}
