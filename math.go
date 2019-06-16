package utils

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

// Abs — Absolute value
func Abs(number float64) float64 {
	return math.Abs(number)
}

// Rand — Generate a random integer
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(math.MaxInt32)

	return n/((math.MaxInt32+1)/(max-min+1)) + min
}

// Round — Rounds a float
func Round(value float64) float64 {
	return math.Floor(value + 0.5)
}

// Floor — Round fractions down
func Floor(value float64) float64 {
	return math.Floor(value)
}

// Ceil — Round fractions up
func Ceil(value float64) float64 {
	return math.Ceil(value)
}

// Pi — Get value of pi
func Pi() float64 {
	return math.Pi
}

// Max — Find highest value
func Max(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}

	max := nums[0]
	for i := 1; i < len(nums); i++ {
		max = math.Max(max, nums[i])
	}

	return max
}

// Min — Find lowest value
func Min(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}

	min := nums[0]
	for i := 1; i < len(nums); i++ {
		min = math.Min(min, nums[i])
	}

	return min
}

// DecBin — Decimal to binary
func DecBin(number int64) string {
	return strconv.FormatInt(number, 2)
}

// BinDec — Binary to decimal
func BinDec(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(i, 10), nil
}

// Hex2Bin — Decodes a hexadecimally encoded binary string
func Hex2Bin(data string) (string, error) {
	i, err := strconv.ParseInt(data, 16, 0)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(i, 2), nil
}

// Bin2Hex — Convert binary data into hexadecimal representation
func Bin2Hex(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(i, 16), nil
}

// DecHex — Decimal to hexadecimal
func DecHex(number int64) string {
	return strconv.FormatInt(number, 16)
}

// HexDec — Hexadecimal to decimal
func HexDec(str string) (int64, error) {
	return strconv.ParseInt(str, 16, 0)
}

// DecOct — Decimal to octal
func DecOct(number int64) string {
	return strconv.FormatInt(number, 8)
}

// OctDec — Octal to decimal
func OctDec(str string) (int64, error) {
	return strconv.ParseInt(str, 8, 0)
}

// BaseConvert — Convert a number between arbitrary bases
func BaseConvert(number string, frombase, tobase int) (string, error) {
	i, err := strconv.ParseInt(number, frombase, 0)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(i, tobase), nil
}

// IsNan — Finds whether a value is not a number
func IsNan(val float64) bool {
	return math.IsNaN(val)
}
