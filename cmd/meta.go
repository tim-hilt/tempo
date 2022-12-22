package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tim-hilt/tempo/util/config"
	"github.com/tim-hilt/tempo/util/logging"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tempo",
	Short: "CLI to communicate with the Tempo-Timesheets API",
	Long:  "CLI to communicate with the Tempo-Timesheets API",
	Run:   root,
}

// TODO: Add interactive command here
func root(cmd *cobra.Command, args []string) {}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, logging.Init, config.Validate)

	rootCmd.PersistentFlags().
		IntP(config.LOGLEVEL_FLAG_VAL, "l", 0, "Logging-level, -1 (trace) to 5 (panic)")
	viper.BindPFlag(
		config.LOGLEVEL_CONFIG_VAL,
		rootCmd.PersistentFlags().Lookup(config.LOGLEVEL_FLAG_VAL),
	)
	viper.SetDefault(config.LOGLEVEL_CONFIG_VAL, 3)

	rootCmd.PersistentFlags().StringP(config.USER_FLAG_VAL, "u", "", "The Jira-User")
	viper.BindPFlag(config.USER_CONFIG_VAL, rootCmd.PersistentFlags().Lookup(config.USER_FLAG_VAL))

	rootCmd.PersistentFlags().
		StringP(config.PASSWORD_FLAG_VAL, "p", "", "The Password for the Jira-User")
	viper.BindPFlag(
		config.PASSWORD_CONFIG_VAL,
		rootCmd.PersistentFlags().Lookup(config.PASSWORD_FLAG_VAL),
	)

	rootCmd.PersistentFlags().StringP(config.JIRA_USER_FLAG_VAL, "juid", "", "The Jira-User-ID")
	viper.BindPFlag(
		config.JIRA_USER_CONFIG_VAL,
		rootCmd.PersistentFlags().Lookup(config.JIRA_USER_FLAG_VAL),
	)

	rootCmd.PersistentFlags().StringP(config.TEMPO_TOKEN_FLAG_VAL, "tt", "", "The Tempo Api Token")
	viper.BindPFlag(
		config.TEMPO_TOKEN_CONFIG_VAL,
		rootCmd.PersistentFlags().Lookup(config.TEMPO_TOKEN_FLAG_VAL),
	)

	rootCmd.PersistentFlags().
		StringP(config.NOTESDIR_FLAG_VAL, "n", "", "The directory of the daily notes")
	viper.BindPFlag(
		config.NOTESDIR_CONFIG_VAL,
		rootCmd.PersistentFlags().Lookup(config.NOTESDIR_FLAG_VAL),
	)
	viper.SetDefault(config.NOTESDIR_CONFIG_VAL, ".")

	rootCmd.PersistentFlags().String(config.HOST_FLAG_VAL, "", "The host of the Jira-instance")
	viper.BindPFlag(config.HOST_CONFIG_VAL, rootCmd.PersistentFlags().Lookup(config.HOST_FLAG_VAL))

	rootCmd.PersistentFlags().
		Bool(config.DEBUG_ENABLED_FLAG_VAL, false, "Whether debugging-information should be enabled for the rest-calls")
	viper.BindPFlag(
		config.DEBUG_ENABLED_CONFIG_VAL,
		rootCmd.PersistentFlags().Lookup(config.DEBUG_ENABLED_FLAG_VAL),
	)

	viper.SetDefault(config.TICKETS_COLUMN_CONFIG_VAL, "Ticket")
	viper.SetDefault(config.COMMENTS_COLUMN_CONFIG_VAL, "Comment")
	viper.SetDefault(config.DURATIONS_COLUMN_CONFIG_VAL, "Duration")
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
