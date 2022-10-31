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

// ticketsForDayCmd represents the ticketsForDay command
var ticketsForDayCmd = &cobra.Command{
	Use:   "tickets-for-day [date]",
	Short: "Print all tickets for given day to the console",
	Long:  "Print all tickets for given day to the console",
	Run:   ticketsForDay,
}

func ticketsForDay(cmd *cobra.Command, args []string) {
	date := parse.ParseArgs(args)
	tempo := tempo.New(flags.User, flags.Password)
	tempo.GetTicketsForDay(date)
}

func init() {
	getCmd.AddCommand(ticketsForDayCmd)
}
