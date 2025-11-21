package commands

import (
	"fmt"

	"github.com/deviantintegral/terminus-golang/pkg/api"
	"github.com/spf13/cobra"
)

var paymentMethodListCmd = &cobra.Command{
	Use:   "payment-method:list",
	Short: "List payment methods",
	Long:  "Display a list of payment methods for your account",
	RunE:  runPaymentMethodList,
}

func init() {
	rootCmd.AddCommand(paymentMethodListCmd)
}

func runPaymentMethodList(_ *cobra.Command, _ []string) error {
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

	methods, err := usersService.ListPaymentMethods(getContext(), sess.UserID)
	if err != nil {
		return fmt.Errorf("failed to list payment methods: %w", err)
	}

	if len(methods) == 0 {
		printMessage("There are no payment methods attached to this account")
		return nil
	}

	return printOutput(methods)
}
