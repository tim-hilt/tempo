package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tim-hilt/tempo/cmd/flags"
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
	homedir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	defaultConfig := filepath.Join(homedir, ".config/tempo/tempo.yaml")

	cobra.OnInitialize(func() { initConfig(defaultConfig) })

	rootCmd.PersistentFlags().StringVar(&flags.Config, "config", defaultConfig, "config file")
	// TODO: Default value is not taken into account
	rootCmd.PersistentFlags().IntVarP(&flags.Loglevel, "loglevel", "l", 3, "Logging-level, -1 (trace) to 5 (panic)")
	rootCmd.PersistentFlags().StringVarP(&flags.User, "user", "u", "", "The Jira-User")
	rootCmd.PersistentFlags().StringVarP(&flags.Password, "password", "p", "", "The Password for the Jira-User")
	rootCmd.PersistentFlags().StringVarP(&flags.NotesDir, "notesdir", "n", ".", "The directory of the daily notes")

	// TODO: Get this to work
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	viper.SetDefault("loglevel", 3)
}

func initConfig(defaultConfig string) {
	if flags.Config != "" {
		viper.SetConfigFile(flags.Config)
	} else {
		viper.AddConfigPath(filepath.Dir(defaultConfig))
		viper.SetConfigName(filepath.Base(defaultConfig))
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
