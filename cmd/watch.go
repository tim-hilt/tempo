package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/tempo"
	"github.com/tim-hilt/tempo/util"
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
	params := util.GetConfigParams()
	tempoClient := tempo.New(params.User, params.Password)

	tempoClient.WatchNotes()
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
