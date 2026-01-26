package commands

import (
	"fmt"
	"time"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/deviantintegral/terminus-golang/pkg/api/models"
	"github.com/spf13/cobra"
)

// formatTimestamp formats a Unix timestamp as a human-readable date string
func formatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return "Never"
	}
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

var backupListCmd = &cobra.Command{
	Use:     "backup:list <site>.<env>",
	Aliases: []string{"backups"},
	Short:   "List backups",
	Long:    "Display a list of backups for an environment",
	Args:    cobra.ExactArgs(1),
	RunE:    runBackupList,
}

var backupCreateCmd = &cobra.Command{
	Use:   "backup:create <site>.<env>",
	Short: "Create a backup",
	Long:  "Create a new backup for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupCreate,
}

var backupGetCmd = &cobra.Command{
	Use:   "backup:get <site>.<env>",
	Short: "Download a backup",
	Long:  "Download a backup to a local file",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupGet,
}

var backupRestoreCmd = &cobra.Command{
	Use:   "backup:restore <site>.<env>",
	Short: "Restore a backup",
	Long:  "Restore an environment from a backup",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupRestore,
}

var backupInfoCmd = &cobra.Command{
	Use:   "backup:info <site>.<env>",
	Short: "Show backup information",
	Long:  "Display detailed information about a specific backup",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupInfo,
}

var backupAutomaticInfoCmd = &cobra.Command{
	Use:   "backup:automatic:info <site>.<env>",
	Short: "Show automatic backup schedule",
	Long:  "Display the automatic backup schedule for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupAutomaticInfo,
}

var backupAutomaticEnableCmd = &cobra.Command{
	Use:   "backup:automatic:enable <site>.<env>",
	Short: "Enable automatic backups",
	Long:  "Enable automatic backup schedule for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupAutomaticEnable,
}

var backupAutomaticDisableCmd = &cobra.Command{
	Use:   "backup:automatic:disable <site>.<env>",
	Short: "Disable automatic backups",
	Long:  "Disable automatic backup schedule for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackupAutomaticDisable,
}

var (
	backupElementFlag string
	backupKeepForFlag int
	backupOutputFlag  string
	backupIDFlag      string
	backupScheduleDay int
)

func init() {
	// Add backup commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(backupListCmd)
	rootCmd.AddCommand(backupCreateCmd)
	rootCmd.AddCommand(backupGetCmd)
	rootCmd.AddCommand(backupRestoreCmd)
	rootCmd.AddCommand(backupInfoCmd)
	rootCmd.AddCommand(backupAutomaticInfoCmd)
	rootCmd.AddCommand(backupAutomaticEnableCmd)
	rootCmd.AddCommand(backupAutomaticDisableCmd)

	// Create flags
	backupCreateCmd.Flags().StringVar(&backupElementFlag, "element", "", "Backup element (code, database, files)")
	backupCreateCmd.Flags().IntVar(&backupKeepForFlag, "keep-for", 0, "Keep backup for N days")

	// Get flags
	backupGetCmd.Flags().StringVar(&backupElementFlag, "element", "code", "Element to download (code, database, files)")
	backupGetCmd.Flags().StringVar(&backupOutputFlag, "output", "", "Output file path")
	backupGetCmd.Flags().StringVar(&backupIDFlag, "backup", "", "Backup ID")

	// Restore flags
	backupRestoreCmd.Flags().StringVar(&backupIDFlag, "backup", "", "Backup ID to restore")
	_ = backupRestoreCmd.MarkFlagRequired("backup")

	// Schedule flags
	backupAutomaticEnableCmd.Flags().IntVar(&backupScheduleDay, "day", 0, "Day of the week for backups (0-6)")

	// Info flags
	backupInfoCmd.Flags().StringVar(&backupElementFlag, "element", "files", "Backup element to show (code, database, files)")
}

func runBackupList(_ *cobra.Command, args []string) error {
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

func runBackupCreate(_ *cobra.Command, args []string) error {
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

func runBackupGet(_ *cobra.Command, args []string) error {
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

func runBackupRestore(_ *cobra.Command, args []string) error {
	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	if !confirm(fmt.Sprintf("Are you sure you want to restore %s.%s from backup %s?", siteID, envID, backupIDFlag)) {
		printMessage("Canceled")
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

func runBackupInfo(_ *cobra.Command, args []string) error {
	siteID, envID, err := parseSiteEnv(args[0])
	if err != nil {
		return err
	}

	backupsService := api.NewBackupsService(cliContext.APIClient)

	// List all backups to find matching one
	backups, err := backupsService.List(getContext(), siteID, envID)
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	// Find the most recent backup of the specified element type
	var targetBackup *models.Backup
	for _, backup := range backups {
		if backup.ArchiveType == backupElementFlag {
			if targetBackup == nil || backup.Timestamp > targetBackup.Timestamp {
				targetBackup = backup
			}
		}
	}

	if targetBackup == nil {
		return fmt.Errorf("no %s backup found for %s.%s", backupElementFlag, siteID, envID)
	}

	// Create a response with backup info and download URL
	type BackupInfo struct {
		ID        string `json:"id"`
		File      string `json:"file"`
		Size      int64  `json:"size"`
		Date      string `json:"date"`
		Expiry    string `json:"expiry"`
		Initiator string `json:"initiator"`
		URL       string `json:"url"`
	}

	// Get download URL
	downloadURL, err := backupsService.GetDownloadURL(getContext(), siteID, envID, targetBackup.ID, backupElementFlag)
	if err != nil {
		// URL might not be available, continue without it
		downloadURL = ""
	}

	info := BackupInfo{
		ID:        targetBackup.ID,
		File:      fmt.Sprintf("%s_%s.tar.gz", targetBackup.ArchiveType, targetBackup.Folder),
		Size:      targetBackup.Size,
		Date:      targetBackup.GetDate().Format("2006-01-02 15:04:05"),
		Expiry:    formatTimestamp(targetBackup.ExpiryTime),
		Initiator: targetBackup.InitiatorEmail,
		URL:       downloadURL,
	}

	return printOutput(info)
}

func runBackupAutomaticInfo(_ *cobra.Command, args []string) error {
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

func runBackupAutomaticEnable(_ *cobra.Command, args []string) error {
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

func runBackupAutomaticDisable(_ *cobra.Command, args []string) error {
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
