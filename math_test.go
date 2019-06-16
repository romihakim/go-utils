package utils

import (
	"testing"
)

func TestMath(t *testing.T) {
	tMax := Max(2, 3.7, 5, 1.1)
	equal(t, float64(5), tMax)

	tMin := Min(2, 3.7, 5, 1.1)
	equal(t, 1.1, tMin)

	tRand := Rand(2, 5)
	rangeValue(t, float64(2), float64(5), float64(tRand))

	tDecbin := DecBin(100)
	equal(t, "1100100", tDecbin)

	tBindec, _ := BinDec(tDecbin)
	equal(t, "100", tBindec)

	tBin2hex, _ := Bin2Hex(tDecbin)
	equal(t, "64", tBin2hex)

	tHexdec, _ := HexDec(tBin2hex)
	equal(t, int64(100), tHexdec)

	tHex2bin, _ := Hex2Bin(tBin2hex)
	equal(t, "1100100", tHex2bin)

	tDecoct := DecOct(tHexdec)
	equal(t, "144", tDecoct)

	tOctdec, _ := OctDec(tDecoct)
	equal(t, int64(100), tOctdec)

	tDechex := DecHex(tHexdec)
	equal(t, "64", tDechex)

	tBaseConvert, _ := BaseConvert("64", 16, 2)
	equal(t, "1100100", tBaseConvert)

	tNumberFormat := NumberFormat(1234567890.777, 2, ".", ",")
	equal(t, "1,234,567,890.78", tNumberFormat)
}
