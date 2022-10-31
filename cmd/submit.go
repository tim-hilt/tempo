package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/flags"
	"github.com/tim-hilt/tempo/cmd/flags/parse"
	"github.com/tim-hilt/tempo/tempo"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit [date]",
	Short: "Submit a daily note to Tempo",
	Long:  `Submit a daily note to Tempo`,
	Run:   submit,
}

func init() {
	rootCmd.AddCommand(submitCmd)
}

func submit(cmd *cobra.Command, args []string) {
	day := parse.ParseArgs(args)

	tempoClient := tempo.New(flags.User, flags.Password)
	tempoClient.SubmitDay(day)
}
