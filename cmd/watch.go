package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/tempo"
	"github.com/tim-hilt/tempo/util/config"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch a given directory and submit changes",
	Long:  "watch a given directory and submit changes",
	Args:  cobra.NoArgs,
	Run:   watch,
}

func watch(cmd *cobra.Command, args []string) {
	user, password := config.GetCredentials()
	tempoClient := tempo.New(user, password)
	tempoClient.WatchNotes()
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
