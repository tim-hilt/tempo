package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/parse"
	"github.com/tim-hilt/tempo/tempo"
	"github.com/tim-hilt/tempo/util/config"
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
	day := parse.ParseDateArg(args)

	user, password := config.GetCredentials()
	tempoClient := tempo.New(user, password)
	tempoClient.SubmitDay(day)
}
