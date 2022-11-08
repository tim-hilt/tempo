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
	Notesdir       string `validate:"required,dir"`
	JiraHost       string `validate:"required,url"`
	Loglevel       int    `validate:"gte=-1,lte=5"`
	DailyWorkhours int    `validate:"gte=0,lte=10"`
}

func GetCredentials() (string, string) {
	return viper.GetString(USER_CONFIG_VAL), viper.GetString(PASSWORD_CONFIG_VAL)
}

func GetHost() string {
	return viper.GetString(HOST_CONFIG_VAL)
}

func GetNotesdir() string {
	notesDir := viper.GetString(NOTESDIR_CONFIG_VAL)
	if strings.HasPrefix(notesDir, "~") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("can't get users homedir")
		}
		notesDir = homedir + notesDir[1:]
	}
	return notesDir
}

func GetWorkhours() int {
	return viper.GetInt(DAILYWORKHOURS_CONFIG_VAL)
}

func GetLoglevel() int {
	return viper.GetInt(LOGLEVEL_CONFIG_VAL)
}

func Validate() {
	user, password := GetCredentials()
	config := configParams{
		User:           user,
		Password:       password,
		JiraHost:       GetHost(),
		Notesdir:       GetNotesdir(),
		Loglevel:       GetLoglevel(),
		DailyWorkhours: GetWorkhours(),
	}
	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		log.Fatal().Err(err).Msg("error when validating config-params")
	}
}
