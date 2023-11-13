package timeconv

import (
	"fmt"
	"time"
)

const dayLayout = "02.01.2006"

const monthLayout = "01.2006"

func TimeToDayString(date time.Time) string {
	return fmt.Sprintf("%02d.%02d.%v", date.Day(), int(date.Month()), date.Year())
}

func TimeToMonthString(date time.Time) string {
	return fmt.Sprintf("%02d.%v", int(date.Month()), date.Year())
}

func StringToDayTime(date string) (time.Time, error) {
	return time.Parse(dayLayout, date)
}

func StringToMonthTime(date string) (time.Time, error) {
	return time.Parse(monthLayout, date)
}
