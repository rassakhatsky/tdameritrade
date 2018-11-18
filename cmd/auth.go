package cmd

import (
	"github.com/rassakhatsky/tdameritrade/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use: "auth",

	Short: "Request authentication token",
	Long:  `Send request to acquire the authentication token from tdameritrade.`,
	Run: func(cmd *cobra.Command, args []string) {
		auth.RequestToken(timeout)
	},
}
