package cli

import (
	amazonring "github.com/jonpmay/go-client-amazon-ring/internal/amazonring"
	"github.com/spf13/cobra"
)

var c, _ = amazonring.NewClient(nil)

var deviceCmd = &cobra.Command{
	Use: "devices",
	Aliases: []string{"device", "d"},
	Short: "Interact with Ring devices",
}

var deviceListCmd = &cobra.Command{
	Use: "list",
	Aliases: []string{"ls"},
	Short: "List Ring devices",
	Run: func (cmd *cobra.Command, args []string) {
		c.Devices()
	},
}

func init() {
	rootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(deviceListCmd)
}