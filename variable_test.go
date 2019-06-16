package utils

import (
	"testing"
)

func TestVariable(t *testing.T) {
	equal(t, true, Empty(""))
	equal(t, true, Empty(0))
	equal(t, true, Empty(0.0))
	equal(t, true, Empty(false))
	equal(t, false, Empty([1]string{}))
	equal(t, true, Empty([]int{}))

	var tIsNumeric bool

	tIsNumeric = IsNumeric("-0xaF")
	equal(t, true, tIsNumeric)

	tIsNumeric = IsNumeric("123456")
	equal(t, true, tIsNumeric)
}
