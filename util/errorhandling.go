package util

import (
	"github.com/rs/zerolog/log"
)

func HandleErr(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}
