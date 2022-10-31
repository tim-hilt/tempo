/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Not implemented yet",
	Long:  "Not implemented yet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command isn't implemented yet.")
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
