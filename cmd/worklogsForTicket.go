package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/tempo"
)

// worklogsForTicketCmd represents the worklogsForTicket command
var worklogsForTicketCmd = &cobra.Command{
	Use:   "worklogs-for-ticket [ticket]",
	Short: "Print all worklogs for the given ticket to the console",
	Long:  "Print all worklogs for the given ticket to the console",
	Run:   worklogsForTicket,
	Args:  cobra.ExactArgs(1),
}

func worklogsForTicket(cmd *cobra.Command, args []string) {
	ticket := args[0]

	tempo.New().GetWorklogsForTicket(ticket)
}

func init() {
	getCmd.AddCommand(worklogsForTicketCmd)
}
