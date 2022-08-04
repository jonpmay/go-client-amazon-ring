package cli

import (
	"context"

	"github.com/spf13/cobra"
	auth "github.com/jonpmay/go-client-amazon-ring/internal/auth"
)

var authCmd = &cobra.Command{
	Use: "auth",
	Short: "Manage the CLI's authentication state",
}

var authLoginCmd = &cobra.Command{
	Use: "login",
	Short: "Authenticate to Ring and retrieve a token",
	Run: func (cmd *cobra.Command, args []string) {
		authRequest := auth.Oauth{}
		auth.Auth(context.Background(), authRequest)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
}