package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tim-hilt/tempo/util"
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
	cobra.OnInitialize(initConfig, logging.Init)

	rootCmd.PersistentFlags().IntP(util.LOGLEVEL_FLAG_VAL, "l", 0, "Logging-level, -1 (trace) to 5 (panic)")
	rootCmd.PersistentFlags().StringP(util.USER_FLAG_VAL, "u", "", "The Jira-User")
	rootCmd.PersistentFlags().StringP(util.PASSWORD_FLAG_VAL, "p", "", "The Password for the Jira-User")
	rootCmd.PersistentFlags().StringP(util.NOTESDIR_FLAG_VAL, "n", "", "The directory of the daily notes")

	viper.BindPFlag(util.LOGLEVEL_CONFIG_VAL, rootCmd.PersistentFlags().Lookup(util.LOGLEVEL_FLAG_VAL))
	viper.BindPFlag(util.USER_CONFIG_VAL, rootCmd.PersistentFlags().Lookup(util.USER_FLAG_VAL))
	viper.BindPFlag(util.PASSWORD_CONFIG_VAL, rootCmd.PersistentFlags().Lookup(util.PASSWORD_FLAG_VAL))
	viper.BindPFlag(util.NOTESDIR_CONFIG_VAL, rootCmd.PersistentFlags().Lookup(util.NOTESDIR_FLAG_VAL))

	viper.SetDefault(util.LOGLEVEL_CONFIG_VAL, 3)
	viper.SetDefault(util.NOTESDIR_CONFIG_VAL, ".")
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
