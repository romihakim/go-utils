package utils

import (
	"testing"
)

func TestDate(t *testing.T) {
	ttimestamp := Time()
	gt(t, float64(ttimestamp), 1522684800)

	tdate := Date("02/01/2006 15:04:05 PM", 1524799394)
	equal(t, "27/04/2018 10:23:14 AM", tdate)

	tstrtotime1, _ := StrToTime("02/01/2006 15:04:05", "02/01/2016 15:04:05")
	equal(t, int64(1451747045), tstrtotime1)

	tstrtotime2, _ := StrToTime("3 04 PM", "8 41 PM")
	equal(t, int64(-62167144740), tstrtotime2)

	equal(t, false, CheckDate(2, 29, 2018))
	equal(t, true, CheckDate(2, 29, 2020))
}
