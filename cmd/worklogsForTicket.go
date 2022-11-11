package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/tempo"
	"github.com/tim-hilt/tempo/util/config"
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

	user, password := config.GetCredentials()
	tempoClient := tempo.New(user, password)
	tempoClient.GetWorklogsForTicket(ticket)
}

func init() {
	getCmd.AddCommand(worklogsForTicketCmd)
}
