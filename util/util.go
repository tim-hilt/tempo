package util

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util/logging"
	"golang.org/x/exp/constraints"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logging.SetLoglevel()
}

const (
	DATE_FORMAT       = "2006-01-02"
	MONTH_FORMAT      = "2006-01"
	SECONDS_IN_MINUTE = 60
	MINUTES_IN_HOUR
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
