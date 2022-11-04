package util

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

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
