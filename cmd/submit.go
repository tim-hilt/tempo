package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/parse"
	"github.com/tim-hilt/tempo/tempo"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit [date]",
	Short: "Submit a daily note to Tempo",
	Long:  `Submit a daily note to Tempo`,
	Run:   submit,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(submitCmd)
}

func submit(cmd *cobra.Command, args []string) {
	date := parse.ParseDateArg(args)

	tempo.New().SubmitDate(date)
}
