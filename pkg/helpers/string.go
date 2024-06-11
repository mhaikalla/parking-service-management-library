package helpers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func MaskingPhone(phone string) string {
	if len(phone) <= 3 {
		return phone
	}
	return strings.Replace(phone[:3]+"***"+phone[len(phone)-3:], "***", "*******", -1)
}

func FormatCommas(num int) string {
	str := fmt.Sprintf("%d", num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != str; {
		n = str
		str = re.ReplaceAllString(str, "$1.$2")
	}
	return str
}

func InsertSpaceNth(s string, n int) string {
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune(' ')
		}
	}
	return buffer.String()
}

func ConvertToUUID(str string) (*uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return nil, err
	}
	return &parsedUUID, nil
}
