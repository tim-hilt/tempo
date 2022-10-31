package parse

import (
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
)

func parseDateArg(dateArg string) string {
	if strings.HasSuffix(dateArg, ".md") {
		pathPieces := strings.Split(dateArg, "/")
		filePieces := strings.Split(pathPieces[len(pathPieces)-1], ".")
		matchDate(filePieces[0])
		return filePieces[0]
	} else if dateArg == "today" {
		return time.Now().Format(util.DATE_FORMAT)
	} else {
		matchDate(dateArg)
		return dateArg
	}
}

func matchDate(date string) {
	if match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}`, date); !match {
		log.Fatal().Msg("date isn't correctly formatted: " + date)
	}
}

func ParseArgs(args []string) string {
	var dateArg string
	if len(args) != 1 {
		log.Fatal().Msg("enter single argument for date")
	} else {
		dateArg = args[0]
	}

	day := parseDateArg(dateArg)

	return day
}
