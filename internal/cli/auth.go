package cli

import (
	"context"
	"fmt"

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
		config := auth.NewConfig(nil)
		authInfo := &auth.AuthInfo{
			Username: "",
			Password: "",
			TwoFactorAuthCode: "",
		}
		token := config.Auth(context.Background(), authInfo)
		fmt.Println(token)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
}
