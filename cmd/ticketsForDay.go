/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/parse"
	"github.com/tim-hilt/tempo/tempo"
	"github.com/tim-hilt/tempo/util"
)

// ticketsForDayCmd represents the ticketsForDay command
var ticketsForDayCmd = &cobra.Command{
	Use:   "tickets-for-day [date]",
	Short: "Print all tickets for given day to the console",
	Long:  "Print all tickets for given day to the console",
	Run:   ticketsForDay,
	Args:  cobra.ExactArgs(1),
}

func ticketsForDay(cmd *cobra.Command, args []string) {
	date := parse.ParseDateArg(args)

	params := util.GetConfigParams()
	tempoClient := tempo.New(params.User, params.Password)

	tempoClient.GetTicketsForDay(date)
}

func init() {
	getCmd.AddCommand(ticketsForDayCmd)
}
