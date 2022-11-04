package util

import (
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
