package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/deviantintegral/terminus-golang/pkg/api/models"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var workflowListCmd = &cobra.Command{
	Use:     "workflow:list <site>",
	Aliases: []string{"workflows"},
	Short:   "List workflows",
	Long:    "Display a list of workflows for a site",
	Args:    cobra.ExactArgs(1),
	RunE:    runWorkflowList,
}

var workflowInfoCmd = &cobra.Command{
	Use:   "workflow:info <site> <workflow-id>",
	Short: "Show workflow information",
	Long:  "Display detailed information about a specific workflow",
	Args:  cobra.ExactArgs(2),
	RunE:  runWorkflowInfo,
}

var workflowWaitCmd = &cobra.Command{
	Use:   "workflow:wait <site> <workflow-id>",
	Short: "Wait for a workflow to complete",
	Long:  "Wait for a workflow to finish and display its status",
	Args:  cobra.ExactArgs(2),
	RunE:  runWorkflowWait,
}

var workflowWatchCmd = &cobra.Command{
	Use:   "workflow:watch <site> <workflow-id>",
	Short: "Watch a workflow",
	Long:  "Watch a workflow and display progress updates",
	Args:  cobra.ExactArgs(2),
	RunE:  runWorkflowWatch,
}

func init() {
	// Add workflow commands directly to rootCmd with colon-separated names
	rootCmd.AddCommand(workflowListCmd)
	rootCmd.AddCommand(workflowInfoCmd)
	rootCmd.AddCommand(workflowWaitCmd)
	rootCmd.AddCommand(workflowWatchCmd)
}

func runWorkflowList(_ *cobra.Command, args []string) error {
	siteID := args[0]
	workflowsService := api.NewWorkflowsService(cliContext.APIClient)

	workflows, err := workflowsService.List(getContext(), siteID)
	if err != nil {
		return fmt.Errorf("failed to list workflows: %w", err)
	}

	return printOutput(workflows)
}

func runWorkflowInfo(_ *cobra.Command, args []string) error {
	siteID := args[0]
	workflowID := args[1]
	workflowsService := api.NewWorkflowsService(cliContext.APIClient)

	workflow, err := workflowsService.Get(getContext(), siteID, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get workflow info: %w", err)
	}

	return printOutput(workflow)
}

func runWorkflowWait(_ *cobra.Command, args []string) error {
	siteID := args[0]
	workflowID := args[1]

	return waitForWorkflow(siteID, workflowID, "Workflow")
}

func runWorkflowWatch(_ *cobra.Command, args []string) error {
	siteID := args[0]
	workflowID := args[1]
	workflowsService := api.NewWorkflowsService(cliContext.APIClient)

	printMessage("Watching workflow %s...", workflowID)

	opts := &api.WatchOptions{
		OnUpdate: func(w *models.Workflow) {
			status := "running"
			if w.IsFinished() {
				if w.IsSuccessful() {
					status = "succeeded"
				} else {
					status = "failed"
				}
			}
			printMessage("[%s] %s - %s", status, w.Type, w.GetMessage())
		},
	}

	if err := workflowsService.Watch(getContext(), siteID, workflowID, opts); err != nil {
		return fmt.Errorf("failed to watch workflow: %w", err)
	}

	// Get final workflow state
	workflow, err := workflowsService.Get(getContext(), siteID, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get final workflow state: %w", err)
	}

	if workflow.IsSuccessful() {
		printMessage("Workflow completed successfully!")
		return nil
	}

	return fmt.Errorf("workflow failed: %s", workflow.GetMessage())
}

// waitForWorkflow waits for a workflow to complete and displays progress
func waitForWorkflow(siteID, workflowID, description string) error {
	workflowsService := api.NewWorkflowsService(cliContext.APIClient)

	// Create progress bar
	var bar *progressbar.ProgressBar
	if !quietFlag {
		bar = progressbar.NewOptions(-1,
			progressbar.OptionSetDescription(description),
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionFullWidth(),
		)
	}

	opts := &api.WaitOptions{
		PollInterval: 3 * time.Second,
		Timeout:      30 * time.Minute,
		OnProgress: func(_ *models.Workflow) {
			if bar != nil {
				_ = bar.Add(1)
			}
		},
	}

	workflow, err := workflowsService.Wait(getContext(), siteID, workflowID, opts)
	if err != nil {
		if bar != nil {
			_ = bar.Finish()
		}
		return fmt.Errorf("workflow wait failed: %w", err)
	}

	if bar != nil {
		_ = bar.Finish()
	}

	if workflow.IsSuccessful() {
		printMessage("%s completed successfully!", description)
		return nil
	}

	return fmt.Errorf("%s failed: %s", description, workflow.GetMessage())
}

// parseSiteEnv parses a site.env string
func parseSiteEnv(input string) (site, env string, err error) {
	parts := strings.SplitN(input, ".", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format: expected 'site.env', got '%s'", input)
	}
	return parts[0], parts[1], nil
}
