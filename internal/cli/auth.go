package cli

import (
	"context"

	auth "github.com/jonpmay/go-client-amazon-ring/internal/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use: "auth",
	Short: "Manage the CLI's authentication state",
}

var authLoginCmd = &cobra.Command{
	Use: "login",
	Short: "Authenticate to Ring and retrieve a token",
	Run: func (cmd *cobra.Command, args []string) {
		c := auth.NewConfig(nil)
		ai := &auth.AuthInfo{
			Username: "",
			Password: "",
			TwoFactorAuthCode: "",
		}
		t := c.Auth(context.Background(), ai)
		creds := auth.NewCredentials(*ai, *t)
		creds.SaveCredentials(creds.EncodeTOMLFile())
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
}
