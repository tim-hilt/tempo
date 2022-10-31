package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tim-hilt/tempo/cmd/flags"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tempo",
	Short: "CLI to communicate with the Tempo-Timesheets API",
	Long:  "CLI to communicate with the Tempo-Timesheets API",
	// TODO: Add watch- or interactive command here
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// TODO: Define config-file at ~/.config/tempo/tempo.yaml
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tempo.yaml)")

	rootCmd.PersistentFlags().IntVarP(&flags.Loglevel, "loglevel", "l", 1, "Logging-level, -1 (trace) to 5 (panic)")
	rootCmd.PersistentFlags().StringVarP(&flags.User, "user", "u", "", "The Jira-User")
	rootCmd.PersistentFlags().StringVarP(&flags.Password, "password", "p", "", "The Password for the Jira-User")

	rootCmd.MarkPersistentFlagRequired("user")
	rootCmd.MarkPersistentFlagRequired("password")
}
