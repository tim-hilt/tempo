package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util/config"
)

func Init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	loglevel := config.GetLoglevel()
	setLoglevel(loglevel)
}

func setLoglevel(loglevel int) {
	if loglevel == -1 {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if loglevel == 0 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if loglevel == 1 {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if loglevel == 2 {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if loglevel == 3 {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if loglevel == 4 {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else if loglevel == 5 {
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	} else {
		log.Fatal().Msg("loglevel has to be between -1 and 5")
	}
}
