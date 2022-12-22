package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/parse"
	"github.com/tim-hilt/tempo/tempo"
)

func init() {
	getCmd.AddCommand(monthlyHoursCmd)
}

// monthlyHoursCmd represents the monthlyHours command
var monthlyHoursCmd = &cobra.Command{
	Use:   "monthly-hours [month]",
	Short: "Get number of monthly hours",
	Long:  "Get number of monthly hours",
	Run:   monthlyHours,
	Args:  cobra.MaximumNArgs(1),
}

func monthlyHours(cmd *cobra.Command, args []string) {
	month := parse.ParseMonthArg(args)

	tempo.New().GetMonthlyHours(month)
}
