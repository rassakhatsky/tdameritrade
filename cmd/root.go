package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	// TODO: add desc
	Use:   "tdgo",
	Short: "tdgo",
	Long:  `tdgo`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hey")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
