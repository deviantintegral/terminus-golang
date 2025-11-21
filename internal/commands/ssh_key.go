package commands

import (
	"fmt"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/spf13/cobra"
)

var sshKeyListCmd = &cobra.Command{
	Use:   "ssh-key:list",
	Short: "List SSH keys",
	Long:  "Display a list of SSH public keys associated with your account",
	RunE:  runSSHKeyList,
}

func init() {
	rootCmd.AddCommand(sshKeyListCmd)
}

func runSSHKeyList(_ *cobra.Command, _ []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	// Load session to get user ID
	sess, err := cliContext.SessionStore.LoadSession()
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}
	if sess == nil || sess.UserID == "" {
		return fmt.Errorf("no user ID in session")
	}

	usersService := api.NewUsersService(cliContext.APIClient)

	keys, err := usersService.ListSSHKeys(getContext(), sess.UserID)
	if err != nil {
		return fmt.Errorf("failed to list SSH keys: %w", err)
	}

	if len(keys) == 0 {
		printMessage("No SSH keys found")
		return nil
	}

	return printOutput(keys)
}
