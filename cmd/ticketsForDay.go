package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/parse"
	"github.com/tim-hilt/tempo/tempo"
	"github.com/tim-hilt/tempo/util/config"
)

// ticketsForDayCmd represents the ticketsForDay command
// TODO: Rename tickets to worklogs
var ticketsForDayCmd = &cobra.Command{
	Use:   "tickets-for-day [date]",
	Short: "Print all tickets for given day to the console",
	Long:  "Print all tickets for given day to the console",
	Run:   ticketsForDay,
	Args:  cobra.ExactArgs(1),
}

func ticketsForDay(cmd *cobra.Command, args []string) {
	date := parse.ParseDateArg(args)

	user, password := config.GetCredentials()
	tempoClient := tempo.New(user, password)
	tempoClient.GetTicketsForDay(date)
}

func init() {
	getCmd.AddCommand(ticketsForDayCmd)
}
