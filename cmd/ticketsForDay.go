package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/parse"
	"github.com/tim-hilt/tempo/tempo"
)

// worklogsForDayCmd represents the ticketsForDay command
var worklogsForDayCmd = &cobra.Command{
	Use:   "worklogs-for-day [date]",
	Short: "Print all worklogs for given day to the console",
	Long:  "Print all worklogs for given day to the console",
	Run:   worklogsForDay,
	Args:  cobra.ExactArgs(1),
}

func worklogsForDay(cmd *cobra.Command, args []string) {
	date := parse.ParseDateArg(args)

	tempo.New().GetWorklogsForDay(date)
}

func init() {
	getCmd.AddCommand(worklogsForDayCmd)
}
