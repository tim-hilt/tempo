/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// overtimeForMonthCmd represents the overtimeForMonth command
var overtimeForMonthCmd = &cobra.Command{
	Use:   "overtime-for-month [month]",
	Short: "Get overtime for a given month",
	Long:  "Get overtime for a given month",
	Run:   overtimeForMonth,
}

func overtimeForMonth(cmd *cobra.Command, args []string) {
	fmt.Println("overtimeForMonth called")
}

func init() {
	getCmd.AddCommand(overtimeForMonthCmd)
}
