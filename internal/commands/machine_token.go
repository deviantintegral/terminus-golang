package commands

import (
	"fmt"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/spf13/cobra"
)

var machineTokenListCmd = &cobra.Command{
	Use:   "machine-token:list",
	Short: "List machine tokens",
	Long:  "Display a list of machine tokens for your account",
	RunE:  runMachineTokenList,
}

func init() {
	rootCmd.AddCommand(machineTokenListCmd)
}

func runMachineTokenList(_ *cobra.Command, _ []string) error {
	// Load session to get user ID
	sess, err := cliContext.SessionStore.LoadSession()
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}
	if sess == nil || sess.UserID == "" {
		return fmt.Errorf("no user ID in session")
	}

	usersService := api.NewUsersService(cliContext.APIClient)

	tokens, err := usersService.ListMachineTokens(getContext(), sess.UserID)
	if err != nil {
		return fmt.Errorf("failed to list machine tokens: %w", err)
	}

	if len(tokens) == 0 {
		printMessage("No machine tokens found")
		return nil
	}

	return printOutput(tokens)
}
