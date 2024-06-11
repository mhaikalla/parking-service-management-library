package helpers

import (
	"regexp"
	"strings"
)

func CreateSlug(s string) string {
	re := regexp.MustCompile(`[^\w\s-]`)
	s = re.ReplaceAllString(strings.ToLower(s), " ")
	re = regexp.MustCompile(`[\s-]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}
