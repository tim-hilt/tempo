package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tim-hilt/tempo/cmd/parse"
	"github.com/tim-hilt/tempo/tempo"
	"github.com/tim-hilt/tempo/util/config"
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

	tempo.New().GetMonthlyOvertime(month)
}

func init() {
	getCmd.AddCommand(overtimeForMonthCmd)

	rootCmd.PersistentFlags().
		IntP(config.DAILYWORKHOURS_FLAG_VAL, "w", 8, "Number of daily working-hours")
	viper.BindPFlag(
		config.DAILYWORKHOURS_CONFIG_VAL,
		rootCmd.PersistentFlags().Lookup(config.DAILYWORKHOURS_FLAG_VAL),
	)
	viper.SetDefault(config.DAILYWORKHOURS_CONFIG_VAL, 8)
}
