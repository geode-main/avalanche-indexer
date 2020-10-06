package util

import (
	"fmt"
	"time"
)

func HourInterval(t time.Time) (time.Time, time.Time) {
	year, month, day := t.Date()

	start := time.Date(year, month, day, t.Hour(), 0, 0, 0, t.Location())
	end := time.Date(year, month, day, t.Hour(), 59, 59, 0, t.Location())

	return start, end
}

func DayInterval(t time.Time) (time.Time, time.Time) {
	year, month, day := t.Date()

	start := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	end := time.Date(year, month, day, 23, 59, 59, 0, t.Location())

	return start, end
}

func TimeBucket(t time.Time, bucket string) time.Time {
	switch bucket {
	case "h":
		start, _ := HourInterval(t)
		return start
	case "d":
		start, _ := DayInterval(t)
		return start
	default:
		panic(fmt.Sprintf("invalid time bucket: %s", bucket))
	}
}
