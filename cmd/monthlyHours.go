package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/tempo"
	"github.com/tim-hilt/tempo/util"
)

func init() {
	getCmd.AddCommand(monthlyHoursCmd)
}

// monthlyHoursCmd represents the monthlyHours command
var monthlyHoursCmd = &cobra.Command{
	Use:   "monthly-hours",
	Short: "Get number of monthly hours",
	Long:  "Get number of monthly hours",
	Args:  cobra.NoArgs,
	Run:   monthlyHours,
}

func monthlyHours(cmd *cobra.Command, args []string) {
	params := util.GetConfigParams()
	tempoClient := tempo.New(params.User, params.Password)
	tempoClient.GetMonthlyHours()
}
