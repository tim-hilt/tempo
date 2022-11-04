package util

import "github.com/spf13/viper"

type configParams struct {
	User           string
	Password       string
	Notesdir       string
	Loglevel       int
	DailyWorkhours int
}

func GetConfigParams() configParams {
	return configParams{
		User:           viper.GetString(USER_CONFIG_VAL),
		Password:       viper.GetString(PASSWORD_CONFIG_VAL),
		Notesdir:       viper.GetString(NOTESDIR_CONFIG_VAL),
		Loglevel:       viper.GetInt(LOGLEVEL_CONFIG_VAL),
		DailyWorkhours: viper.GetInt(DAILYWORKHOURS_CONFIG_VAL),
	}
}
