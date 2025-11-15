package commands

import (
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/pantheon-systems/terminus-go/pkg/api/models"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:     "backup",
	Aliases: []string{"backups"},
	Short:   "Backup management commands",
	Long:    "Manage site backups",
}

var backupListCmd = &cobra.Command{
	Use:     "list <site>.<env>",
	Aliases: []string{"ls"},
	Short:   "List backups",
	Long:    "Display a list of backups for an environment",
	Args:    cobra.ExactArgs(1),
	RunE:    runBackupList,
}

var backupCreateCmd = &cobra.Command{
	Use:   "create <site>.<env>",
	Short: "Create a backup",
	Long:  "Create a new backup for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupCreate,
}

var backupGetCmd = &cobra.Command{
	Use:   "get <site>.<env>",
	Short: "Download a backup",
	Long:  "Download a backup to a local file",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupGet,
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore <site>.<env>",
	Short: "Restore a backup",
	Long:  "Restore an environment from a backup",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupRestore,
}

var backupAutomaticCmd = &cobra.Command{
	Use:   "automatic",
	Short: "Automatic backup management",
	Long:  "Manage automatic backup schedules",
}

var backupAutomaticInfoCmd = &cobra.Command{
	Use:   "info <site>.<env>",
	Short: "Show automatic backup schedule",
	Long:  "Display the automatic backup schedule for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupAutomaticInfo,
}

var backupAutomaticEnableCmd = &cobra.Command{
	Use:   "enable <site>.<env>",
	Short: "Enable automatic backups",
	Long:  "Enable automatic backup schedule for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupAutomaticEnable,
}

var backupAutomaticDisableCmd = &cobra.Command{
	Use:   "disable <site>.<env>",
	Short: "Disable automatic backups",
	Long:  "Disable automatic backup schedule for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupAutomaticDisable,
}

var (
	backupElementFlag  string
	backupKeepForFlag  int
	backupOutputFlag   string
	backupIDFlag       string
	backupScheduleDay  int
)

func init() {
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupCreateCmd)
	backupCmd.AddCommand(backupGetCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	backupCmd.AddCommand(backupAutomaticCmd)

	backupAutomaticCmd.AddCommand(backupAutomaticInfoCmd)
	backupAutomaticCmd.AddCommand(backupAutomaticEnableCmd)
	backupAutomaticCmd.AddCommand(backupAutomaticDisableCmd)

	// Create flags
	backupCreateCmd.Flags().StringVar(&backupElementFlag, "element", "", "Backup element (code, database, files)")
	backupCreateCmd.Flags().IntVar(&backupKeepForFlag, "keep-for", 0, "Keep backup for N days")

	// Get flags
	backupGetCmd.Flags().StringVar(&backupElementFlag, "element", "code", "Element to download (code, database, files)")
	backupGetCmd.Flags().StringVar(&backupOutputFlag, "output", "", "Output file path")
	backupGetCmd.Flags().StringVar(&backupIDFlag, "backup", "", "Backup ID")

	// Restore flags
	backupRestoreCmd.Flags().StringVar(&backupIDFlag, "backup", "", "Backup ID to restore")
	backupRestoreCmd.MarkFlagRequired("backup")

	// Schedule flags
	backupAutomaticEnableCmd.Flags().IntVar(&backupScheduleDay, "day", 0, "Day of the week for backups (0-6)")
}

func runBackupList(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	backupsService := api.NewBackupsService(cliContext.APIClient)

	backups, err := backupsService.List(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	return printOutput(backups)
}

func runBackupCreate(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	backupsService := api.NewBackupsService(cliContext.APIClient)

	var workflow *models.Workflow

	if backupElementFlag != "" {
		printMessage("Creating %s backup for %s.%s...", backupElementFlag, siteID, envID)
		workflow, err = backupsService.CreateElement(getContext(), siteID, envID, backupElementFlag)
	} else {
		printMessage("Creating backup for %s.%s...", siteID, envID)
		req := &api.CreateBackupRequest{
			KeepFor: backupKeepForFlag,
		}
		workflow, err = backupsService.Create(getContext(), siteID, envID, req)
	}

	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Creating backup")
}

func runBackupGet(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	if backupIDFlag == "" {
		return fmt.Errorf("--backup flag is required")
	}

	backupsService := api.NewBackupsService(cliContext.APIClient)

	// Determine output path
	outputPath := backupOutputFlag
	if outputPath == "" {
		outputPath = fmt.Sprintf("%s-%s-%s-%s.tar.gz", siteID, envID, backupIDFlag, backupElementFlag)
	}

	printMessage("Downloading %s backup to %s...", backupElementFlag, outputPath)

	if err := backupsService.Download(getContext(), siteID, envID, backupIDFlag, backupElementFlag, outputPath); err != nil {
		return fmt.Errorf("failed to download backup: %w", err)
	}

	printMessage("Backup downloaded successfully!")

	return nil
}

func runBackupRestore(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	if !confirm(fmt.Sprintf("Are you sure you want to restore %s.%s from backup %s?", siteID, envID, backupIDFlag)) {
		printMessage("Cancelled")
		return nil
	}

	backupsService := api.NewBackupsService(cliContext.APIClient)

	printMessage("Restoring backup %s for %s.%s...", backupIDFlag, siteID, envID)

	workflow, err := backupsService.Restore(getContext(), siteID, envID, backupIDFlag)
	if err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	return waitForWorkflow(siteID, workflow.ID, "Restoring backup")
}

func runBackupAutomaticInfo(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	backupsService := api.NewBackupsService(cliContext.APIClient)

	schedule, err := backupsService.GetSchedule(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to get backup schedule: %w", err)
	}

	return printOutput(schedule)
}

func runBackupAutomaticEnable(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	backupsService := api.NewBackupsService(cliContext.APIClient)

	printMessage("Enabling automatic backups for %s.%s...", siteID, envID)

	if err := backupsService.SetSchedule(getContext(), siteID, envID, true, backupScheduleDay); err != nil {
		return fmt.Errorf("failed to enable automatic backups: %w", err)
	}

	printMessage("Automatic backups enabled!")

	return nil
}

func runBackupAutomaticDisable(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	backupsService := api.NewBackupsService(cliContext.APIClient)

	printMessage("Disabling automatic backups for %s.%s...", siteID, envID)

	if err := backupsService.SetSchedule(getContext(), siteID, envID, false, 0); err != nil {
		return fmt.Errorf("failed to disable automatic backups: %w", err)
	}

	printMessage("Automatic backups disabled!")

	return nil
}
