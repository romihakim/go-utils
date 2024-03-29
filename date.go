package utils

import (
	"time"
)

// Time — Return current Unix timestamp
func Time() int64 {
	return time.Now().Unix()
}

// StrToTime — Parse about any English textual datetime description into a Unix timestamp
// StrToTime("02/01/2006 15:04:05", "02/01/2016 15:04:05") == 1451747045
// StrToTime("3 04 PM", "8 41 PM") == -62167144740
func StrToTime(format, strtime string) (int64, error) {
	t, err := time.Parse(format, strtime)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

// Date — Format a local time/date
// Date("02/01/2006 15:04:05 PM", 1524799394)
func Date(format string, timestamp int64) string {
	return time.Unix(timestamp, 0).Format(format)
}

// CheckDate — Validate a Gregorian date
func CheckDate(month, day, year int) bool {
	if month < 1 || month > 12 || day < 1 || day > 31 || year < 1 || year > 32767 {
		return false
	}

	switch month {
	case 4, 6, 9, 11:
		if day > 30 {
			return false
		}
	case 2:
		// leap year
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			if day > 29 {
				return false
			}
		} else if day > 28 {
			return false
		}
	}

	return true
}

// Sleep — Delay execution
func Sleep(t int64) {
	time.Sleep(time.Duration(t) * time.Second)
}

// Usleep — Delay execution in microseconds
func Usleep(t int64) {
	time.Sleep(time.Duration(t) * time.Microsecond)
}
