package commands

import (
	"fmt"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/spf13/cobra"
)

var envListCmd = &cobra.Command{
	Use:     "env:list <site>",
	Aliases: []string{"environment"},
	Short:   "List environments",
	Long:    "Display a list of all environments for a site",
	Args:    cobra.ExactArgs(1),
	RunE:    runEnvList,
}

var envInfoCmd = &cobra.Command{
	Use:   "env:info <site>.<env>",
	Short: "Show environment information",
	Long:  "Display detailed information about a specific environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvInfo,
}

var envClearCacheCmd = &cobra.Command{
	Use:   "env:clear-cache <site>.<env>",
	Short: "Clear environment cache",
	Long:  "Clear the cache for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvClearCache,
}

var envDeployCmd = &cobra.Command{
	Use:   "env:deploy <site>.<env>",
	Short: "Deploy code to an environment",
	Long:  "Deploy code from one environment to another",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvDeploy,
}

var envCloneContentCmd = &cobra.Command{
	Use:   "env:clone-content <site>.<env>",
	Short: "Clone content between environments",
	Long:  "Clone database and/or files from one environment to another",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvCloneContent,
}

var envCommitCmd = &cobra.Command{
	Use:   "env:commit <site>.<env>",
	Short: "Commit changes",
	Long:  "Commit changes in SFTP mode",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvCommit,
}

var envWipeCmd = &cobra.Command{
	Use:   "env:wipe <site>.<env>",
	Short: "Wipe environment",
	Long:  "Wipe content from an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvWipe,
}

var envConnectionSetCmd = &cobra.Command{
	Use:   "env:connection:set <site>.<env> <mode>",
	Short: "Set connection mode",
	Long:  "Set the connection mode (git or sftp)",
	Args:  cobra.ExactArgs(2),
	RunE:  runEnvConnectionSet,
}

var envMetricsCmd = &cobra.Command{
	Use:   "env:metrics <site>[.<env>]",
	Short: "Display environment metrics",
	Long:  "Display pages served and unique visit metrics for a site environment. Use <site>.<env> for environment-specific metrics, or <site> for aggregated site metrics.",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvMetrics,
}

var (
	envUpdateDBFlag   bool
	envNoteFlag       string
	envClearCacheFlag bool
	envFromEnvFlag    string
	envDatabaseFlag   bool
	envFilesFlag      bool
	envCommitMsgFlag  string
	envMetricsPeriod  string
	envMetricsDatapts string
)

func init() {
	// Add env commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(envListCmd)
	rootCmd.AddCommand(envInfoCmd)
	rootCmd.AddCommand(envClearCacheCmd)
	rootCmd.AddCommand(envDeployCmd)
	rootCmd.AddCommand(envCloneContentCmd)
	rootCmd.AddCommand(envCommitCmd)
	rootCmd.AddCommand(envWipeCmd)
	rootCmd.AddCommand(envConnectionSetCmd)

	// Deploy flags
	envDeployCmd.Flags().BoolVar(&envUpdateDBFlag, "updatedb", false, "Run database updates after deploy")
	envDeployCmd.Flags().StringVar(&envNoteFlag, "note", "", "Deploy note/annotation")
	envDeployCmd.Flags().BoolVar(&envClearCacheFlag, "cc", true, "Clear cache after deploy")

	// Clone content flags
	envCloneContentCmd.Flags().StringVar(&envFromEnvFlag, "from-env", "", "Source environment")
	envCloneContentCmd.Flags().BoolVar(&envDatabaseFlag, "db", true, "Clone database")
	envCloneContentCmd.Flags().BoolVar(&envFilesFlag, "files", true, "Clone files")
	_ = envCloneContentCmd.MarkFlagRequired("from-env")

	// Commit flags
	envCommitCmd.Flags().StringVarP(&envCommitMsgFlag, "message", "m", "", "Commit message")
	_ = envCommitCmd.MarkFlagRequired("message")

	// Metrics command
	rootCmd.AddCommand(envMetricsCmd)

	// Metrics flags
	envMetricsCmd.Flags().StringVar(&envMetricsPeriod, "period", "day", "Time period for metrics: day, week, or month")
	envMetricsCmd.Flags().StringVar(&envMetricsDatapts, "datapoints", "auto", "Number of data points to return, or 'auto' for intelligent default")
}

func runEnvList(_ *cobra.Command, args []string) error {
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

func runEnvInfo(_ *cobra.Command, args []string) error {
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

func runEnvClearCache(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	if !confirm(fmt.Sprintf("Are you sure you want to clear the cache for %s.%s?", siteID, envID)) {
		printMessage("Canceled")
		return nil
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	printMessage("Clearing cache for %s.%s...", siteID, envID)

	workflow, err := envsService.ClearCache(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Clearing cache")
}

func runEnvDeploy(_ *cobra.Command, args []string) error {
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

func runEnvCloneContent(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	// Build confirmation message based on what's being cloned
	var cloneTypes []string
	if envDatabaseFlag {
		cloneTypes = append(cloneTypes, "database")
	}
	if envFilesFlag {
		cloneTypes = append(cloneTypes, "files")
	}
	cloneTypeStr := "content"
	if len(cloneTypes) > 0 {
		cloneTypeStr = fmt.Sprintf("%v", cloneTypes)
	}

	if !confirm(fmt.Sprintf("Are you sure you want to clone %s from %s to %s.%s? This will overwrite existing content.", cloneTypeStr, envFromEnvFlag, siteID, envID)) {
		printMessage("Canceled")
		return nil
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

func runEnvCommit(_ *cobra.Command, args []string) error {
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

func runEnvWipe(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	if !confirm(fmt.Sprintf("Are you sure you want to wipe %s.%s? This cannot be undone.", siteID, envID)) {
		printMessage("Canceled")
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

func runEnvConnectionSet(_ *cobra.Command, args []string) error {
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

func runEnvMetrics(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	// Determine duration string from period and datapoints
	duration, err := buildMetricsDuration(envMetricsPeriod, envMetricsDatapts)
	if err != nil {
		return err
	}

	// Parse the input - could be "site" or "site.env"
	input := args[0]
	var siteID, envID string

	// Check if input contains a period (site.env format)
	if idx := findEnvSeparator(input); idx != -1 {
		siteID, envID, err = parseSiteEnv(input)
		if err != nil {
			return err
		}
	} else {
		// Site-only format - get aggregated metrics
		siteID = input
		envID = ""
	}

	envsService := api.NewEnvironmentsService(cliContext.APIClient)

	metrics, err := envsService.GetMetrics(getContext(), siteID, envID, duration)
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}

	return printOutput(metrics)
}

// buildMetricsDuration builds the duration string for the metrics API
func buildMetricsDuration(period, datapoints string) (string, error) {
	// Map period to short form and default datapoints
	var shortForm string
	var defaultDatapoints int

	switch period {
	case "day":
		shortForm = "d"
		defaultDatapoints = 28
	case "week":
		shortForm = "w"
		defaultDatapoints = 12
	case "month":
		shortForm = "m"
		defaultDatapoints = 12
	default:
		return "", fmt.Errorf("invalid period: %s (must be 'day', 'week', or 'month')", period)
	}

	// Determine number of datapoints
	var numDatapoints int
	if datapoints == "auto" {
		numDatapoints = defaultDatapoints
	} else {
		n, err := fmt.Sscanf(datapoints, "%d", &numDatapoints)
		if err != nil || n != 1 || numDatapoints < 1 {
			return "", fmt.Errorf("invalid datapoints: %s (must be a positive number or 'auto')", datapoints)
		}
	}

	return fmt.Sprintf("%d%s", numDatapoints, shortForm), nil
}

// findEnvSeparator finds the position of the environment separator (period)
// Returns -1 if no separator is found
func findEnvSeparator(input string) int {
	for i, c := range input {
		if c == '.' {
			return i
		}
	}
	return -1
}
