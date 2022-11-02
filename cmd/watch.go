/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/flags"
	"github.com/tim-hilt/tempo/tempo"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch a given directory and submit changes",
	Long:  "watch a given directory and submit changes",
	Args:  cobra.NoArgs,
	Run:   watch,
}

func watch(cmd *cobra.Command, args []string) {
	tempoClient := tempo.New(flags.JiraUser, flags.Password)
	tempoClient.WatchNotes()
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
