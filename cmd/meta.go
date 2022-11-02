package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tim-hilt/tempo/util/logging"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tempo",
	Short: "CLI to communicate with the Tempo-Timesheets API",
	Long:  "CLI to communicate with the Tempo-Timesheets API",
	Run:   root,
}

// TODO: Add watch- or interactive command here
func root(cmd *cobra.Command, args []string) {}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, logging.SetLoglevel)

	rootCmd.PersistentFlags().IntP("loglevel", "l", 0, "Logging-level, -1 (trace) to 5 (panic)")
	rootCmd.PersistentFlags().StringP("jirauser", "u", "", "The Jira-User")
	rootCmd.PersistentFlags().StringP("password", "p", "", "The Password for the Jira-User")
	rootCmd.PersistentFlags().StringP("notesdir", "n", "", "The directory of the daily notes")

	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("jiraUser", rootCmd.PersistentFlags().Lookup("jirauser"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("notesDir", rootCmd.PersistentFlags().Lookup("notesdir"))

	viper.SetDefault("loglevel", 3)
	viper.SetDefault("notesDir", ".")
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(filepath.Join(home, ".config/tempo/"))
	viper.SetConfigType("yaml")
	viper.SetConfigName("tempo")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	cobra.CheckErr(err)
}
