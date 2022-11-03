/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tim-hilt/tempo/cmd/parse"
	"github.com/tim-hilt/tempo/tempo"
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

	user := viper.GetString("jiraUser")
	password := viper.GetString("password")
	tempoClient := tempo.New(user, password)

	tempoClient.GetMonthlyOvertime(month)
}

func init() {
	getCmd.AddCommand(overtimeForMonthCmd)

	rootCmd.PersistentFlags().IntP("workhours", "w", 8, "Number of daily working-hours")
	viper.BindPFlag("dailyWorkhours", rootCmd.PersistentFlags().Lookup("workhours"))
	viper.SetDefault("dailyWorkhours", 8)
}
