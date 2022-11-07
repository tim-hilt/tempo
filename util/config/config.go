package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type configParams struct {
	User           string `validate:"required"`
	Password       string `validate:"required"`
	Notesdir       string `validate:"required"`
	JiraHost       string `validate:"required"`
	Loglevel       int
	DailyWorkhours int
}

func GetConfigParams() configParams {
	notesDir := formatNotesDir(viper.GetString(NOTESDIR_CONFIG_VAL))
	config := configParams{
		User:           viper.GetString(USER_CONFIG_VAL),
		Password:       viper.GetString(PASSWORD_CONFIG_VAL),
		JiraHost:       viper.GetString(HOST_CONFIG_VAL),
		Notesdir:       notesDir,
		Loglevel:       viper.GetInt(LOGLEVEL_CONFIG_VAL),
		DailyWorkhours: viper.GetInt(DAILYWORKHOURS_CONFIG_VAL),
	}
	return config
}

func Validate() {
	notesDir := formatNotesDir(viper.GetString(NOTESDIR_CONFIG_VAL))
	config := configParams{
		User:           viper.GetString(USER_CONFIG_VAL),
		Password:       viper.GetString(PASSWORD_CONFIG_VAL),
		JiraHost:       viper.GetString(HOST_CONFIG_VAL),
		Notesdir:       notesDir,
		Loglevel:       viper.GetInt(LOGLEVEL_CONFIG_VAL),
		DailyWorkhours: viper.GetInt(DAILYWORKHOURS_CONFIG_VAL),
	}
	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		log.Fatal().Err(err).Msg("error when validating config-params")
	}
}

func formatNotesDir(notesDir string) string {
	if strings.HasPrefix(notesDir, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("error when formating notesDir")
		}
		return homeDir + notesDir[1:]
	}
	return notesDir
}
