package parse

import (
	"regexp"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
)

func ParseDateArg(args []string) string {
	dateArg := args[0]
	if dateArg == "today" {
		return time.Now().Format(util.DATE_FORMAT)
	} else {
		validateDate(dateArg)
		return dateArg
	}
}

func validateDate(date string) {
	if match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}`, date); !match {
		log.Fatal().Str("date", date).Msg("date isn't correctly formatted")
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
		log.Fatal().Str("date", date).Msg("date isn't correctly formatted")
	}
}
