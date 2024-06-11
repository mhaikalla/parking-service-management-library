package helpers

import (
	"time"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func DiffSecondsFromNow(timeToCompare time.Time) int {
	return DiffSeconds(timeToCompare, time.Now())
}

func DiffSeconds(fromTime time.Time, toTime time.Time) int {
	//fmt.Println("Now = " + fromTime.Format("2006.01.02 15:04:05"))
	datetimeNow := time.Date(toTime.Year(), toTime.Month(), toTime.Day(), toTime.Hour(), toTime.Minute(), toTime.Second(), 0, time.UTC)
	//fmt.Println("Date Time Now = " + datetimeNow.Format("2006.01.02 15:04:05"))
	diffSeconds := fromTime.Sub(datetimeNow).Seconds()
	diffSecondInt := Abs(int(diffSeconds))
	//fmt.Println("Diff Seconds = " + strconv.Itoa(diffSecondInt))
	return diffSecondInt
}
