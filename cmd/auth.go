package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gotd", // TODO: find cool name
	Short: "gotd",
	Long:  ``, // TODO: Add long desc
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("failed to execute the command")
		os.Exit(1)
	}
}
