package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/pantheon-systems/terminus-go/pkg/config"
	"github.com/pantheon-systems/terminus-go/pkg/output"
	"github.com/pantheon-systems/terminus-go/pkg/session"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	formatFlag   string
	fieldsFlag   []string
	yesFlag      bool
	quietFlag    bool
	verboseCount int
)

// CLIContext holds shared context for all commands
type CLIContext struct {
	Config       *config.Config
	SessionStore *session.Store
	APIClient    *api.Client
	Output       *output.Options
}

var cliContext *CLIContext

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "terminus",
	Short: "Terminus Go - Pantheon command line interface",
	Long: `Terminus Go is a command line interface for managing Pantheon sites.

It provides tools to manage sites, environments, workflows, backups, and more
on the Pantheon platform.`,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		return initCLIContext()
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "table", "Output format (table, json, yaml, csv, list)")
	rootCmd.PersistentFlags().StringSliceVar(&fieldsFlag, "fields", nil, "Fields to display (comma-separated)")
	rootCmd.PersistentFlags().BoolVarP(&yesFlag, "yes", "y", false, "Answer yes to all prompts")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress output")
	rootCmd.PersistentFlags().CountVarP(&verboseCount, "verbose", "v", "Verbose output (-v, -vv, or -vvv for increasing verbosity)")

	// Note: All commands are now added directly to rootCmd in their respective files using colon-separated names:
	// - auth commands (auth:login, auth:logout, auth:whoami) in auth.go
	// - site commands (site:list, site:info, site:org:list, etc.) in site.go
	// - env commands (env:list, env:info, etc.) in env.go
	// - workflow commands (workflow:list, workflow:info, etc.) in workflow.go
	// - backup commands (backup:list, backup:create, etc.) in backup.go
	// - org commands (org:list, org:info, etc.) in org.go
	// - domain commands (domain:list, domain:add, etc.) in domain.go
	// - multidev commands (multidev:create, multidev:delete, multidev:list, etc.) in multidev.go
	// - connection commands (connection:info, connection:set) in connection.go
	// - lock commands (lock:info, lock:enable, lock:disable) in lock.go
	// - plan commands (plan:info, plan:list) in plan.go
	// - upstream commands (upstream:info, upstream:list, upstream:updates:list) in upstream.go
	// - self commands (self:info) in self.go
	// - art commands (art, art:list) in art.go
	// - redis commands (redis:enable, redis:disable) in redis.go
	// - branch commands (branch:list) in branch.go
	// - machine-token commands (machine-token:list) in machine_token.go
	// - payment-method commands (payment-method:list) in payment_method.go
	// - ssh-key commands (ssh-key:list) in ssh_key.go
	// - tag commands (tag:list) in tag.go
}

// initCLIContext initializes the CLI context
func initCLIContext() error {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create session store
	sessionStore := session.NewStore(cfg.CacheDir)

	// Create API client with optional logger
	clientOpts := []api.ClientOption{
		api.WithBaseURL(cfg.GetBaseURL()),
	}

	// Add logger if verbose mode is enabled
	if verboseCount > 0 {
		logger := api.NewLogger(api.VerbosityLevel(verboseCount))
		clientOpts = append(clientOpts, api.WithLogger(logger))
	}

	apiClient := api.NewClient(clientOpts...)

	// Try to load existing session
	sess, err := sessionStore.LoadSession()
	if err == nil && sess != nil {
		apiClient.SetToken(sess.SessionToken)
	}

	// Create output options
	outputOpts := &output.Options{
		Format: output.Format(formatFlag),
		Fields: fieldsFlag,
		Writer: os.Stdout,
	}

	cliContext = &CLIContext{
		Config:       cfg,
		SessionStore: sessionStore,
		APIClient:    apiClient,
		Output:       outputOpts,
	}

	return nil
}

// requireAuth ensures the user is authenticated
func requireAuth() error {
	if cliContext.APIClient == nil {
		return fmt.Errorf("API client not initialized")
	}

	// Check if we have a session
	sess, err := cliContext.SessionStore.LoadSession()
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}

	if sess == nil {
		return fmt.Errorf("not authenticated. Please run 'terminus auth:login' first")
	}

	return nil
}

// confirm prompts the user for confirmation
func confirm(message string) bool {
	if yesFlag {
		return true
	}

	fmt.Printf("%s [y/N]: ", message)
	var response string
	_, _ = fmt.Scanln(&response)

	return response == "y" || response == "Y" || response == "yes"
}

// printOutput prints data using the configured output format
func printOutput(data interface{}) error {
	if quietFlag {
		return nil
	}

	return output.Print(data, cliContext.Output)
}

// printMessage prints a message to stdout
func printMessage(format string, args ...interface{}) {
	if !quietFlag {
		fmt.Printf(format+"\n", args...)
	}
}

// printError prints an error message to stderr
func printError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}

// getContext returns a context for API calls
func getContext() context.Context {
	return context.Background()
}
