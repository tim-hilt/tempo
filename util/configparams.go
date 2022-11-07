package util

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type configParams struct {
	User           string
	Password       string
	Notesdir       string
	Loglevel       int
	DailyWorkhours int
}

func GetConfigParams() configParams {
	notesDir := formatNotesDir(viper.GetString(NOTESDIR_CONFIG_VAL))
	return configParams{
		User:           viper.GetString(USER_CONFIG_VAL),
		Password:       viper.GetString(PASSWORD_CONFIG_VAL),
		Notesdir:       notesDir,
		Loglevel:       viper.GetInt(LOGLEVEL_CONFIG_VAL),
		DailyWorkhours: viper.GetInt(DAILYWORKHOURS_CONFIG_VAL),
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
