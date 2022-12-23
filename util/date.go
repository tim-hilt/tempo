package util

import (
	"regexp"
)

func IsFullDate(date string) bool {
	match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}`, date)
	return match
}

func IsYearMonth(date string) bool {
	match, _ := regexp.MatchString(`\d{4}-\d{2}`, date)
	return match
}
