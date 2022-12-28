package util

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func Divmod[T constraints.Integer](numerator, denominator T) (T, T) {
	quotient := numerator / denominator
	remainder := numerator % denominator
	return quotient, remainder
}

func Max[T constraints.Ordered](args ...T) (max T) {
	for _, val := range args {
		if val > max {
			max = val
		}
	}
	return
}

func Min[T constraints.Ordered](args ...T) (min T) {
	for _, val := range args {
		if val < min {
			min = val
		}
	}
	return
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SlicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func CalcDurationSeconds(duration string) (int, error) {
	if duration == "" {
		return 0, nil
	}

	minutesHours := strings.Split(duration, ":")
	hours, err := strconv.Atoi(minutesHours[0])

	if err != nil {
		return -1, err
	}

	minutes, err := strconv.Atoi(minutesHours[1])

	if err != nil {
		return -1, err
	}

	return (hours*MINUTES_IN_HOUR + minutes) * SECONDS_IN_MINUTE, nil
}

func FromPreviousMonths(day time.Time) bool {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	return day.Before(firstOfMonth)
}
