package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tim-hilt/tempo/cmd/parse"
	"github.com/tim-hilt/tempo/tempo"
	"github.com/tim-hilt/tempo/util"
)

// overtimeForMonthCmd represents the overtimeForMonth command
var overtimeForMonthCmd = &cobra.Command{
	Use:   "overtime-for-month [month]",
	Short: "Get overtime for a given month",
	Long:  "Get overtime for a given month",
	Run:   overtimeForMonth,
	Args:  cobra.MaximumNArgs(1),
}

func overtimeForMonth(cmd *cobra.Command, args []string) {
	month := parse.ParseMonthArg(args)

	params := util.GetConfigParams()
	tempoClient := tempo.New(params.User, params.Password)
	tempoClient.GetMonthlyOvertime(month)
}

func init() {
	getCmd.AddCommand(overtimeForMonthCmd)

	rootCmd.PersistentFlags().IntP(util.DAILYWORKHOURS_FLAG_VAL, "w", 8, "Number of daily working-hours")
	viper.BindPFlag(util.DAILYWORKHOURS_CONFIG_VAL, rootCmd.PersistentFlags().Lookup(util.DAILYWORKHOURS_FLAG_VAL))
	viper.SetDefault(util.DAILYWORKHOURS_CONFIG_VAL, 8)
}
