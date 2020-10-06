package util

import (
	"strconv"
	"time"
)

func ParseInt64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}

func ParseUnixTime(val string) (time.Time, error) {
	secs, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(secs, 0), nil
}

func ParseFloat32(val string) (float64, error) {
	return strconv.ParseFloat(val, 32)
}
