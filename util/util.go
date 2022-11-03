package util

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/exp/constraints"
)

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

func HandleErr(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}

func handleErronousHttpStatus(status int) {
	if status < 200 || status > 299 {
		log.Fatal().Msg("http-status was " + fmt.Sprint(status) + " instead of 200")
	}
}

func HandleResponse(status int, err error, msg string) {
	HandleErr(err, msg)
	handleErronousHttpStatus(status)
}

func Divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func Max[T constraints.Ordered](args ...T) (max T) {
	for _, val := range args {
		if val > max {
			max = val
		}
	}
	return
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
