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

	user := viper.GetString("jiraUser")
	password := viper.GetString("password")
	tempoClient := tempo.New(user, password)

	tempoClient.GetTicketsForDay(date)
}

func init() {
	getCmd.AddCommand(ticketsForDayCmd)
}
