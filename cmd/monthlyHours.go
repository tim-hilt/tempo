package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Args:  cobra.NoArgs,
	Run:   monthlyHours,
}

func monthlyHours(cmd *cobra.Command, args []string) {
	user := viper.GetString("jiraUser")
	password := viper.GetString("password")
	tempoClient := tempo.New(user, password)
	tempoClient.GetMonthlyHours()
}
