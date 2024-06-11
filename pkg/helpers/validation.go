package helpers

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsEmailValid checks if the email provided passes the required structure and length.
func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

// IsPhoneNumberValid checks if the msisdn provided passes the required structure and length.
func IsPhoneNumberValid(phoneNumber string, local string) bool {
	newPhoneNumber := strings.Replace(phoneNumber[0:2], "08", "628", -1) + phoneNumber[2:]
	lengthNumber := len(newPhoneNumber)

	if !(lengthNumber >= 10 && lengthNumber <= 13) {
		return false
	}

	pattern := `^\+?[1-9]\d{1,14}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(phoneNumber)
}
