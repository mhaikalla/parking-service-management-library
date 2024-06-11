package helpers

import (
	"strings"
	"unicode"
)

// StringContains ..
func StringContains(slices []string, comparison string) bool {
	for _, a := range slices {
		if a == comparison {
			return true
		}
	}

	return false
}

// IntContains ..
func IntContains(slices []int, comparison int) bool {
	for _, a := range slices {
		if a == comparison {
			return true
		}
	}

	return false
}

// Float64Contains ..
func Float64Contains(slices []float64, comparison int) bool {
	for _, a := range slices {
		if int(a) == comparison {
			return true
		}
	}

	return false
}

// StringInSlice ..
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// URLContains ...
func URLContains(comparison string, slices []string) bool {
	for _, a := range slices {
		if strings.Contains(comparison, a) {
			return true
		}
	}
	return false
}

// IsInt ...
func IsInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
