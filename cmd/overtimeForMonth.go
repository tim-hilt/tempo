/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/flags"
	"github.com/tim-hilt/tempo/cmd/flags/parse"
	"github.com/tim-hilt/tempo/tempo"
)

// overtimeForMonthCmd represents the overtimeForMonth command
var overtimeForMonthCmd = &cobra.Command{
	Use:   "overtime-for-month [month]",
	Short: "Get overtime for a given month",
	Long:  "Get overtime for a given month",
	Run:   overtimeForMonth,
}

func overtimeForMonth(cmd *cobra.Command, args []string) {
	month := parse.ParseMonthArg(args)

	tempoClient := tempo.New(flags.User, flags.Password)
	tempoClient.GetMonthlyOvertime(month)
}

func init() {
	getCmd.AddCommand(overtimeForMonthCmd)

	rootCmd.PersistentFlags().IntVarP(&flags.DailyWorkhours, "workhours", "w", 8, "Number of daily working-hours")
}
