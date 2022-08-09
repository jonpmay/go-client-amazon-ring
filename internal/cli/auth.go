package cli

import (
	"fmt"

	"github.com/jonpmay/go-client-amazon-ring/internal/amazonring"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use: "auth",
	Short: "Manage the CLI's authentication state",
}

var authLoginCmd = &cobra.Command{
	Use: "login",
	Short: "Authenticate to Ring using username and password to retrieve a token",
	Run: func (cmd *cobra.Command, args []string) {
		c, err := amazonring.NewClientWithOptions(nil, amazonring.SetBaseURL(amazonring.AuthURL))
		if err != nil {
			panic(err)
		}
		t := c.PasswordGrant()
		fmt.Println(t)
	},
}

var authLoginRefreshCmd = &cobra.Command{
	Use: "refresh",
	Short: "Authenticate to Ring using an existing refresh token and retrieve a new token",
	Run: func (cmd *cobra.Command, args []string) {
		// TODO
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
}
