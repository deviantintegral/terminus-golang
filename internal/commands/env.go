package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:     "env",
	Aliases: []string{"environment"},
	Short:   "Environment management commands",
	Long:    "Manage site environments",
}

var envListCmd = &cobra.Command{
	Use:     "list <site>",
	Aliases: []string{"ls"},
	Short:   "List environments",
	Long:    "Display a list of all environments for a site",
	Args:    cobra.ExactArgs(1),
	RunE:    runEnvList,
}

var envInfoCmd = &cobra.Command{
	Use:   "info <site>.<env>",
	Short: "Show environment information",
	Long:  "Display detailed information about a specific environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvInfo,
}

var envClearCacheCmd = &cobra.Command{
	Use:   "clear-cache <site>.<env>",
	Short: "Clear environment cache",
	Long:  "Clear the cache for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvClearCache,
}

var envDeployCmd = &cobra.Command{
	Use:   "deploy <site>.<env>",
	Short: "Deploy code to an environment",
	Long:  "Deploy code from one environment to another",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvDeploy,
}

var envCloneContentCmd = &cobra.Command{
	Use:   "clone-content <site>.<env>",
	Short: "Clone content between environments",
	Long:  "Clone database and/or files from one environment to another",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvCloneContent,
}

var envCommitCmd = &cobra.Command{
	Use:   "commit <site>.<env>",
	Short: "Commit changes",
	Long:  "Commit changes in SFTP mode",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvCommit,
}

var envWipeCmd = &cobra.Command{
	Use:   "wipe <site>.<env>",
	Short: "Wipe environment",
	Long:  "Wipe content from an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvWipe,
}

var envConnectionSetCmd = &cobra.Command{
	Use:   "set <site>.<env> <mode>",
	Short: "Set connection mode",
	Long:  "Set the connection mode (git or sftp)",
	Args:  cobra.ExactArgs(2),
	RunE:  runEnvConnectionSet,
}

var envConnectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "Connection mode management",
	Long:  "Manage environment connection mode",
}

var (
	envUpdateDBFlag   bool
	envNoteFlag       string
	envClearCacheFlag bool
	envFromEnvFlag    string
	envDatabaseFlag   bool
	envFilesFlag      bool
	envCommitMsgFlag  string
)

func init() {
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envInfoCmd)
	envCmd.AddCommand(envClearCacheCmd)
	envCmd.AddCommand(envDeployCmd)
	envCmd.AddCommand(envCloneContentCmd)
	envCmd.AddCommand(envCommitCmd)
	envCmd.AddCommand(envWipeCmd)
	envCmd.AddCommand(envConnectionCmd)

	envConnectionCmd.AddCommand(envConnectionSetCmd)

	// Deploy flags
	envDeployCmd.Flags().BoolVar(&envUpdateDBFlag, "updatedb", false, "Run database updates after deploy")
	envDeployCmd.Flags().StringVar(&envNoteFlag, "note", "", "Deploy note/annotation")
	envDeployCmd.Flags().BoolVar(&envClearCacheFlag, "cc", true, "Clear cache after deploy")

	// Clone content flags
	envCloneContentCmd.Flags().StringVar(&envFromEnvFlag, "from-env", "", "Source environment")
	envCloneContentCmd.Flags().BoolVar(&envDatabaseFlag, "db", true, "Clone database")
	envCloneContentCmd.Flags().BoolVar(&envFilesFlag, "files", true, "Clone files")
	envCloneContentCmd.MarkFlagRequired("from-env")

	// Commit flags
	envCommitCmd.Flags().StringVarP(&envCommitMsgFlag, "message", "m", "", "Commit message")
	envCommitCmd.MarkFlagRequired("message")
}

func runEnvList(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID := args[0]
	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	envs, err := envsService.List(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to list environments: %w", err)
	}

	return printOutput(envs)
}

func runEnvInfo(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	env, err := envsService.Get(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to get environment info: %w", err)
	}

	return printOutput(env)
}

func runEnvClearCache(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	printMessage("Clearing cache for %s.%s...", siteID, envID)

	workflow, err := envsService.ClearCache(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Clearing cache")
}

func runEnvDeploy(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	req := &api.DeployRequest{
		UpdateDB:   envUpdateDBFlag,
		Note:       envNoteFlag,
		ClearCache: envClearCacheFlag,
	}

	printMessage("Deploying to %s.%s...", siteID, envID)

	workflow, err := envsService.Deploy(getContext(), siteID, envID, req)
	if err != nil {
		return fmt.Errorf("failed to deploy: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Deploying")
}

func runEnvCloneContent(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	req := &api.CloneContentRequest{
		FromEnvironment: envFromEnvFlag,
		Database:        envDatabaseFlag,
		Files:           envFilesFlag,
	}

	printMessage("Cloning content to %s.%s from %s...", siteID, envID, envFromEnvFlag)

	workflow, err := envsService.CloneContent(getContext(), siteID, envID, req)
	if err != nil {
		return fmt.Errorf("failed to clone content: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Cloning content")
}

func runEnvCommit(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	req := &api.CommitRequest{
		Message: envCommitMsgFlag,
	}

	printMessage("Committing changes in %s.%s...", siteID, envID)

	workflow, err := envsService.Commit(getContext(), siteID, envID, req)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Committing")
}

func runEnvWipe(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	if !confirm(fmt.Sprintf("Are you sure you want to wipe %s.%s? This cannot be undone.", siteID, envID)) {
		printMessage("Cancelled")
		return nil
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	printMessage("Wiping %s.%s...", siteID, envID)

	workflow, err := envsService.Wipe(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to wipe environment: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Wiping environment")
}

func runEnvConnectionSet(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	mode := args[1]
	if mode != "git" && mode != "sftp" {
		return fmt.Errorf("invalid mode: %s (must be 'git' or 'sftp')", mode)
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	printMessage("Setting connection mode to %s for %s.%s...", mode, siteID, envID)

	workflow, err := envsService.ChangeConnectionMode(getContext(), siteID, envID, mode)
	if err != nil {
		return fmt.Errorf("failed to change connection mode: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Changing connection mode")
}
