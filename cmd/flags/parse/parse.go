package parse

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/logging"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logging.SetLoglevel()
}

func ParseDateArg(args []string) string {
	dateArg := args[0]
	if strings.HasSuffix(dateArg, ".md") {
		pathPieces := strings.Split(dateArg, "/")
		filePieces := strings.Split(pathPieces[len(pathPieces)-1], ".")
		validateDate(filePieces[0])
		return filePieces[0]
	} else if dateArg == "today" {
		return time.Now().Format(util.DATE_FORMAT)
	} else {
		validateDate(dateArg)
		return dateArg
	}
}

func validateDate(date string) {
	if match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}`, date); !match {
		log.Fatal().Msg("date isn't correctly formatted: " + date)
	}
}

func ParseMonthArg(args []string) string {
	var month string
	if len(args) == 0 {
		month = time.Now().Format(util.MONTH_FORMAT)
	} else if len(args) == 1 {
		month = args[0]
	}

	validateMonth(month)

	return month
}

func validateMonth(date string) {
	if match, _ := regexp.MatchString(`\d{4}-\d{2}`, date); !match {
		log.Fatal().Msg("date isn't correctly formatted: " + date)
	}
}
