package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/deviantintegral/terminus-golang/pkg/api/models"
	"github.com/spf13/cobra"
)

var aliasesCmd = &cobra.Command{
	Use:     "aliases",
	Aliases: []string{"drush:aliases"},
	Short:   "Generate Drush aliases for Pantheon sites",
	Long: `Generates Pantheon Drush aliases for sites on which the currently logged-in user is on the team.

Note that Drush 9+ does not read alias files from global locations. You must set valid
alias locations in your drush.yml file. Refer to https://docs.pantheon.io/guides/drush/drush-aliases
for more information.`,
	RunE: runAliases,
}

var (
	aliasesPrintFlag    bool
	aliasesLocationFlag string
	aliasesTypeFlag     string
	aliasesBaseFlag     string
)

func init() {
	rootCmd.AddCommand(aliasesCmd)

	aliasesCmd.Flags().BoolVar(&aliasesPrintFlag, "print", false, "Display aliases to stdout instead of writing files")
	aliasesCmd.Flags().StringVar(&aliasesLocationFlag, "location", "", "Custom file path for Drush 8 aliases (--type must be 'php')")
	aliasesCmd.Flags().StringVar(&aliasesTypeFlag, "type", "all", "Output format: 'php' (Drush 8), 'yml' (Drush 9+), or 'all'")
	aliasesCmd.Flags().StringVar(&aliasesBaseFlag, "base", "~/.drush", "Base directory for alias files")
}

func runAliases(_ *cobra.Command, _ []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	// Validate flags
	if aliasesLocationFlag != "" && aliasesTypeFlag != "php" {
		return fmt.Errorf("--location flag can only be used with --type=php")
	}

	// Load session to get user ID
	sess, err := cliContext.SessionStore.LoadSession()
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}
	if sess == nil || sess.UserID == "" {
		return fmt.Errorf("no user ID in session")
	}

	// Fetch site information
	printMessage("Fetching site information to build Drush aliases...")

	sitesService := api.NewSitesService(cliContext.APIClient)
	sites, err := sitesService.List(getContext(), sess.UserID)
	if err != nil {
		return fmt.Errorf("failed to fetch sites: %w", err)
	}

	if len(sites) == 0 {
		printMessage("No sites found.")
		return nil
	}

	printMessage("%d sites found.", len(sites))

	// Generate and emit aliases based on type flag
	if aliasesPrintFlag {
		// Print to stdout
		return printAliases(sites)
	}

	// Write to files
	return writeAliasFiles(sites)
}

func printAliases(sites []*models.Site) error {
	switch aliasesTypeFlag {
	case "php":
		printMessage("Displaying Drush 8 alias file contents:")
		fmt.Println(generateDrush8Aliases(sites))
	case "yml":
		printMessage("Displaying Drush 9+ alias files:")
		for _, site := range sites {
			fmt.Printf("# %s.site.yml\n", site.Name)
			fmt.Println(generateDrush9Alias(site))
			fmt.Println()
		}
	case "all":
		printMessage("Displaying Drush 8 alias file contents:")
		fmt.Println(generateDrush8Aliases(sites))
		fmt.Println()
		printMessage("Displaying Drush 9+ alias files:")
		for _, site := range sites {
			fmt.Printf("# %s.site.yml\n", site.Name)
			fmt.Println(generateDrush9Alias(site))
			fmt.Println()
		}
	default:
		return fmt.Errorf("invalid --type value: %s (must be 'php', 'yml', or 'all')", aliasesTypeFlag)
	}

	return nil
}

func writeAliasFiles(sites []*models.Site) error {
	baseDir := expandHomePath(aliasesBaseFlag)

	switch aliasesTypeFlag {
	case "php":
		return writeDrush8Aliases(sites, baseDir)
	case "yml":
		return writeDrush9Aliases(sites, baseDir)
	case "all":
		if err := writeDrush8Aliases(sites, baseDir); err != nil {
			return err
		}
		return writeDrush9Aliases(sites, baseDir)
	default:
		return fmt.Errorf("invalid --type value: %s (must be 'php', 'yml', or 'all')", aliasesTypeFlag)
	}
}

func writeDrush8Aliases(sites []*models.Site, baseDir string) error {
	var filePath string
	if aliasesLocationFlag != "" {
		filePath = expandHomePath(aliasesLocationFlag)
	} else {
		filePath = filepath.Join(baseDir, "pantheon.aliases.drushrc.php")
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Generate aliases content
	content := generateDrush8Aliases(sites)

	// Write to file
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to write Drush 8 aliases: %w", err)
	}

	printMessage("Writing Drush 8 alias file to %s", shortenHomePath(filePath))

	return nil
}

func writeDrush9Aliases(sites []*models.Site, baseDir string) error {
	aliasDir := filepath.Join(baseDir, "sites", "pantheon")

	// Ensure directory exists
	if err := os.MkdirAll(aliasDir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", aliasDir, err)
	}

	// Write one file per site
	for _, site := range sites {
		filePath := filepath.Join(aliasDir, fmt.Sprintf("%s.site.yml", site.Name))
		content := generateDrush9Alias(site)

		if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
			return fmt.Errorf("failed to write Drush 9 alias for %s: %w", site.Name, err)
		}
	}

	printMessage("Writing Drush 9 alias files to %s", shortenHomePath(aliasDir))

	return nil
}

func generateDrush8Aliases(sites []*models.Site) string {
	var sb strings.Builder

	// Write header
	sb.WriteString("<?php\n")
	sb.WriteString("  /**\n")
	sb.WriteString("   * Pantheon drush alias file, to be placed in your ~/.drush directory or the aliases\n")
	sb.WriteString("   * directory of your local Drush home. Once it's in place, clear drush cache:\n")
	sb.WriteString("   *\n")
	sb.WriteString("   * drush cc drush\n")
	sb.WriteString("   *\n")
	sb.WriteString("   * To see all your available aliases:\n")
	sb.WriteString("   *\n")
	sb.WriteString("   * drush sa\n")
	sb.WriteString("   *\n")
	sb.WriteString("   * See http://helpdesk.getpantheon.com/customer/portal/articles/411388 for details.\n")
	sb.WriteString("   */\n\n")

	// Write alias for each site
	for _, site := range sites {
		sb.WriteString(fmt.Sprintf("  $aliases['%s.*'] = array(\n", site.Name))
		sb.WriteString(fmt.Sprintf("    'uri' => '${env-name}-%s.pantheonsite.io',\n", site.Name))
		sb.WriteString(fmt.Sprintf("    'remote-host' => 'appserver.${env-name}.%s.drush.in',\n", site.ID))
		sb.WriteString(fmt.Sprintf("    'remote-user' => '${env-name}.%s',\n", site.ID))
		sb.WriteString("    'ssh-options' => '-p 2222 -o \"AddressFamily inet\"',\n")
		sb.WriteString("    'path-aliases' => array(\n")
		sb.WriteString("      '%files' => 'files',\n")
		sb.WriteString("    ),\n")
		sb.WriteString("  );\n\n")
	}

	return sb.String()
}

func generateDrush9Alias(site *models.Site) string {
	var sb strings.Builder

	sb.WriteString("'*':\n")
	sb.WriteString(fmt.Sprintf("  host: appserver.${env-name}.%s.drush.in\n", site.ID))
	sb.WriteString("  paths:\n")
	sb.WriteString("    files: files\n")
	sb.WriteString(fmt.Sprintf("  uri: ${env-name}-%s.pantheonsite.io\n", site.Name))
	sb.WriteString(fmt.Sprintf("  user: ${env-name}.%s\n", site.ID))
	sb.WriteString("  ssh:\n")
	sb.WriteString("    options: '-p 2222 -o \"AddressFamily inet\"'\n")
	sb.WriteString("    tty: false\n")
	sb.WriteString("\n") // Extra newline at end to match PHP terminus

	return sb.String()
}

// expandHomePath expands ~ to home directory
func expandHomePath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// shortenHomePath converts /home/user to ~ for display
func shortenHomePath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}
