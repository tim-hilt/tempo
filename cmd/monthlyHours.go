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
	Run: func(cmd *cobra.Command, args []string) {
		tempoClient := tempo.New(flags.User, flags.Password)
		tempoClient.GetMonthlyHours() // TODO: Somehow the number seems too high! Maybe I have to look into the output again.
	},
}
