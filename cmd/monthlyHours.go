package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/flags"
	"github.com/tim-hilt/tempo/tempo"
)

func init() {
	getCmd.AddCommand(monthlyHoursCmd)
}

// monthlyHoursCmd represents the monthlyHours command
var monthlyHoursCmd = &cobra.Command{
	Use:   "monthly-hours",
	Short: "Get number of monthly hours",
	Long:  "Get number of monthly hours",
	Args:  cobra.ExactArgs(0),
	Run:   monthlyHours,
}

func monthlyHours(cmd *cobra.Command, args []string) {
	tempoClient := tempo.New(flags.User, flags.Password)
	tempoClient.GetMonthlyHours()
}
