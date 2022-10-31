/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ticketsForDayCmd represents the ticketsForDay command
var ticketsForDayCmd = &cobra.Command{
	Use:   "ticketsForDay",
	Short: "Print all tickets for given day to the console",
	Long:  "Print all tickets for given day to the console",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Standalone call not functional")
	},
}

func init() {
	getCmd.AddCommand(ticketsForDayCmd)
}
