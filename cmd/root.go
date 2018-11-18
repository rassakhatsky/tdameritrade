package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// flags are here
var (
	timeout int
)

var rootCmd = &cobra.Command{
	// TODO: add desc
	Use:   "tdgo",
	Short: "tdgo",
	Long:  `tdgo`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please use help to check all available commands:")
		fmt.Println("> tdgo help")
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 90, "Default timeout")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(authCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
