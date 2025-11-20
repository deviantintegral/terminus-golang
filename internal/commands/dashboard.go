package commands

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pantheon-systems/terminus-go/pkg/api"
	"github.com/spf13/cobra"
)

var dashboardViewCmd = &cobra.Command{
	Use:   "dashboard:view [site[.env]]",
	Short: "Open the Pantheon Dashboard in a browser",
	Long: `Display the URL for the Pantheon Dashboard or open the Dashboard in a browser.

Usage examples:
  dashboard:view                  Open your account dashboard
  dashboard:view --print          Print your account dashboard URL
  dashboard:view <site>           Open the specified site's dashboard
  dashboard:view <site>.<env>     Open a specific environment's dashboard`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDashboardView,
}

var printURLFlag bool

func init() {
	rootCmd.AddCommand(dashboardViewCmd)
	dashboardViewCmd.Flags().BoolVar(&printURLFlag, "print", false, "Print URL instead of opening browser")
}

func runDashboardView(_ *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	var dashboardURL string
	var err error

	if len(args) == 0 {
		// No arguments: show user's account dashboard
		dashboardURL, err = getUserDashboardURL()
		if err != nil {
			return err
		}
	} else {
		// Parse argument to determine if it's site or site.env
		input := args[0]
		if strings.Contains(input, ".") {
			// site.env format
			dashboardURL, err = getEnvironmentDashboardURL(input)
			if err != nil {
				return err
			}
		} else {
			// site format
			dashboardURL, err = getSiteDashboardURL(input)
			if err != nil {
				return err
			}
		}
	}

	if printURLFlag {
		// Just print the URL
		_, _ = fmt.Println(dashboardURL)
		return nil
	}

	// Open in browser
	if err := openBrowser(dashboardURL); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	printMessage("Opening %s", dashboardURL)
	return nil
}

func getUserDashboardURL() (string, error) {
	// Load session to get user ID
	sess, err := cliContext.SessionStore.LoadSession()
	if err != nil {
		return "", fmt.Errorf("failed to load session: %w", err)
	}
	if sess == nil || sess.UserID == "" {
		return "", fmt.Errorf("no user ID in session")
	}

	return fmt.Sprintf("https://dashboard.pantheon.io/users/%s", sess.UserID), nil
}

func getSiteDashboardURL(siteID string) (string, error) {
	// Verify site exists and get its ID
	sitesService := api.NewSitesService(cliContext.APIClient)
	site, err := sitesService.Get(getContext(), siteID)
	if err != nil {
		return "", fmt.Errorf("failed to get site info: %w", err)
	}

	return fmt.Sprintf("https://dashboard.pantheon.io/sites/%s", site.ID), nil
}

func getEnvironmentDashboardURL(siteEnv string) (string, error) {
	siteID, envID, err := parseSiteEnv(siteEnv)
	if err != nil {
		return "", err
	}

	// Verify site exists and get its ID
	sitesService := api.NewSitesService(cliContext.APIClient)
	site, err := sitesService.Get(getContext(), siteID)
	if err != nil {
		return "", fmt.Errorf("failed to get site info: %w", err)
	}

	// Verify environment exists
	envsService := api.NewEnvironmentsService(cliContext.APIClient)
	_, err = envsService.Get(getContext(), site.ID, envID)
	if err != nil {
		return "", fmt.Errorf("failed to get environment info: %w", err)
	}

	return fmt.Sprintf("https://dashboard.pantheon.io/sites/%s#%s", site.ID, envID), nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default: // linux, freebsd, openbsd, netbsd
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start() //nolint:gosec // User-controlled URL is intentional
}
