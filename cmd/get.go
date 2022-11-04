package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Parent command for getting several other metrics",
	Long:  "Parent command for getting several other metrics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Standalone call not functional")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
