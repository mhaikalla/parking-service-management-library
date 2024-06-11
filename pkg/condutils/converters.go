package condutils

import (
	"fmt"
	"regexp"
)

// IDRLayoutFormatter format amount to representatif of string IDR layout.
func IDRLayoutFormatter(amount int64, sep string) string {
	re := regexp.MustCompile(`^.+([0-9]{3})$`)
	resStr := ""

	for {
		if amount > 999 {
			strAm := fmt.Sprintf("%v", amount)
			resStr = sep + re.FindStringSubmatch(strAm)[1] + resStr
			amount = int64(amount / 1000)
		}

		if amount < 1000 {
			resStr = fmt.Sprintf("%v", amount) + resStr
			break
		}

	}
	return resStr
}
