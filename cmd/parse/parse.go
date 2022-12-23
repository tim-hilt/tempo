package parse

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
)

func ParseDateArg(args []string) string {
	dateArg := args[0]

	if dateArg == "today" {
		return time.Now().Format(util.DATE_FORMAT)
	} else if dateArg == "month" {
		return time.Now().Format(util.MONTH_FORMAT)
	}

	if !util.IsFullDate(dateArg) {
		log.Fatal().
			Str("date", dateArg).
			Msg("date isn't correctly formatted; expected " + util.DATE_FORMAT)
	}

	return dateArg
}

func ParseMonthArg(args []string) string {
	var month string

	if len(args) == 0 {
		month = time.Now().Format(util.MONTH_FORMAT)
	} else if len(args) == 1 {
		month = args[0]
	}

	if !util.IsYearMonth(month) {
		log.Fatal().
			Str("date", month).
			Msg("date isn't correctly formatted; expected " + util.MONTH_FORMAT)
	}

	return month
}
