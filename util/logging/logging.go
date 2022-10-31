package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/cmd/flags"
)

func SetLoglevel() {
	if flags.Loglevel == -1 {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if flags.Loglevel == 0 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if flags.Loglevel == 1 {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if flags.Loglevel == 2 {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if flags.Loglevel == 3 {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if flags.Loglevel == 4 {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else if flags.Loglevel == 5 {
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	} else {
		log.Fatal().Msg("loglevel has to be between -1 and 5")
	}
}
