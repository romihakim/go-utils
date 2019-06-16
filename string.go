package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"html"
	"io/ioutil"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// StrPos — Find the position of the first occurrence of a substring in a string
func StrPos(haystack, needle string, offset int) int {
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		offset += length
	}

	pos := strings.Index(haystack[offset:], needle)
	if pos == -1 {
		return -1
	}

	return pos + offset
}

// StrIpos — Find the position of the first occurrence of a case-insensitive substring in a string
func StrIpos(haystack, needle string, offset int) int {
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	haystack = haystack[offset:]
	if offset < 0 {
		offset += length
	}

	pos := strings.Index(strings.ToLower(haystack), strings.ToLower(needle))
	if pos == -1 {
		return -1
	}

	return pos + offset
}

// StrRpos — Find the position of the last occurrence of a substring in a string
func StrRpos(haystack, needle string, offset int) int {
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}

	pos = strings.LastIndex(haystack, needle)
	if offset > 0 && pos != -1 {
		pos += offset
	}

	return pos
}

// StrRipos — Find the position of the last occurrence of a case-insensitive substring in a string
func StrRipos(haystack, needle string, offset int) int {
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}

	pos = strings.LastIndex(strings.ToLower(haystack), strings.ToLower(needle))
	if offset > 0 && pos != -1 {
		pos += offset
	}

	return pos
}

// StrReplace — Replace all occurrences of the search string with the replacement string
func StrReplace(search, replace, subject string, count int) string {
	return strings.Replace(subject, search, replace, count)
}

// StrToUpper — Make a string uppercase
func StrToUpper(str string) string {
	return strings.ToUpper(str)
}

// StrToLower — Make a string lowercase
func StrToLower(str string) string {
	return strings.ToLower(str)
}

// UcFirst — Make a string's first character uppercase
func UcFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}

	return ""
}

// LcFirst — Make a string's first character lowercase
func LcFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}

	return ""
}

// UcWords — Uppercase the first character of each word in a string
func UcWords(str string) string {
	return strings.Title(str)
}

// SubStr — Return part of a string
func SubStr(str string, start uint, length int) string {
	if start < 0 || length < -1 {
		return str
	}

	switch {
	case length == -1:
		return str[start:]
	case length == 0:
		return ""
	}

	end := int(start) + length
	if end > len(str) {
		end = len(str)
	}

	return str[start:end]
}

// StrRev — Reverse a string
func StrRev(str string) string {
	runes := []rune(str)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// ParseStr — Parses the string into variables
// f1=m&f2=n -> map[f1:m f2:n]
// f[a]=m&f[b]=n -> map[f:map[a:m b:n]]
// f[a][a]=m&f[a][b]=n -> map[f:map[a:map[a:m b:n]]]
// f[]=m&f[]=n -> map[f:[m n]]
// f[a][]=m&f[a][]=n -> map[f:map[a:[m n]]]
// f[][]=m&f[][]=n -> map[f:[map[]]] // Currently does not support nested slice.
// f=m&f[a]=n -> error // This is not the same as PHP.
// a .[[b=c -> map[a___[b:c]
func ParseStr(encodedString string, result map[string]interface{}) error {
	// build nested map.
	var build func(map[string]interface{}, []string, interface{}) error

	build = func(result map[string]interface{}, keys []string, value interface{}) error {
		length := len(keys)
		// trim ',"
		key := strings.Trim(keys[0], "'\"")
		if length == 1 {
			result[key] = value
			return nil
		}

		// The end is slice. like f[], f[a][]
		if keys[1] == "" && length == 2 {
			// todo nested slice
			if key == "" {
				return nil
			}

			val, ok := result[key]
			if !ok {
				result[key] = []interface{}{value}
				return nil
			}

			children, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("expected type '[]interface{}' for key '%s', but got '%T'", key, val)
			}

			result[key] = append(children, value)
			return nil
		}

		// The end is slice + map. like f[][a]
		if keys[1] == "" && length > 2 && keys[2] != "" {
			val, ok := result[key]
			if !ok {
				result[key] = []interface{}{}
				val = result[key]
			}

			children, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("expected type '[]interface{}' for key '%s', but got '%T'", key, val)
			}

			if l := len(children); l > 0 {
				if child, ok := children[l-1].(map[string]interface{}); ok {
					if _, ok := child[keys[2]]; !ok {
						build(child, keys[2:], value)
						return nil
					}
				}
			}

			child := map[string]interface{}{}
			build(child, keys[2:], value)
			result[key] = append(children, child)

			return nil
		}

		// map. like f[a], f[a][b]
		val, ok := result[key]
		if !ok {
			result[key] = map[string]interface{}{}
			val = result[key]
		}

		children, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected type 'map[string]interface{}' for key '%s', but got '%T'", key, val)
		}

		if err := build(children, keys[1:], value); err != nil {
			return err
		}

		return nil
	}

	// split encodedString.
	parts := strings.Split(encodedString, "&")
	for _, part := range parts {
		pos := strings.Index(part, "=")
		if pos <= 0 {
			continue
		}

		key, err := url.QueryUnescape(part[:pos])
		if err != nil {
			return err
		}

		for key[0] == ' ' {
			key = key[1:]
		}

		if key == "" || key[0] == '[' {
			continue
		}

		value, err := url.QueryUnescape(part[pos+1:])
		if err != nil {
			return err
		}

		// split into multiple keys
		var keys []string
		left := 0
		for i, k := range key {
			if k == '[' && left == 0 {
				left = i
			} else if k == ']' {
				if left > 0 {
					if len(keys) == 0 {
						keys = append(keys, key[:left])
					}

					keys = append(keys, key[left+1:i])
					left = 0

					if i+1 < len(key) && key[i+1] != '[' {
						break
					}
				}
			}
		}

		if len(keys) == 0 {
			keys = append(keys, key)
		}

		// first key
		first := ""
		for i, chr := range keys[0] {
			if chr == ' ' || chr == '.' || chr == '[' {
				first += "_"
			} else {
				first += string(chr)
			}

			if chr == '[' {
				first += keys[0][i+1:]
				break
			}
		}

		keys[0] = first

		// build nested map
		if err := build(result, keys, value); err != nil {
			return err
		}
	}

	return nil
}

// NumberFormat — Format a number with grouped thousands
// decimals: Sets the number of decimal points.
// decPoint: Sets the separator for the decimal point.
// thousandsSep: Sets the thousands separator.
func NumberFormat(number float64, decimals uint, decPoint, thousandsSep string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}

	dec := int(decimals)
	// Will round off
	str := fmt.Sprintf("%."+strconv.Itoa(dec)+"F", number)
	prefix, suffix := "", ""
	if dec > 0 {
		prefix = str[:len(str)-(dec+1)]
		suffix = str[len(str)-dec:]
	} else {
		prefix = str
	}

	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	// thousands sep num
	c := (l1 - 1) / 3
	tmp := make([]byte, l2*c+l1)
	pos := len(tmp) - 1

	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}

		tmp[pos] = prefix[i]
	}

	s := string(tmp)
	if dec > 0 {
		s += decPoint + suffix
	}

	if neg {
		s = "-" + s
	}

	return s
}

// ChunkSplit — Split a string into smaller chunks
func ChunkSplit(body string, chunklen uint, end string) string {
	if end == "" {
		end = "\r\n"
	}

	runes, erunes := []rune(body), []rune(end)
	l := uint(len(runes))
	if l <= 1 || l < chunklen {
		return body + end
	}

	ns := make([]rune, 0, len(runes)+len(erunes))
	var i uint
	for i = 0; i < l; i += chunklen {
		if i+chunklen > l {
			ns = append(ns, runes[i:]...)
		} else {
			ns = append(ns, runes[i:i+chunklen]...)
		}
		ns = append(ns, erunes...)
	}

	return string(ns)
}

// StrWordCount — Return information about words used in a string
func StrWordCount(str string) []string {
	return strings.Fields(str)
}

// WordWrap — Wraps a string to a given number of characters
func WordWrap(str string, width uint, br string) string {
	if br == "" {
		br = "\n"
	}

	init := make([]byte, 0, len(str))
	buf := bytes.NewBuffer(init)
	var current uint
	var wordbuf, spacebuf bytes.Buffer

	for _, char := range str {
		if char == '\n' {
			if wordbuf.Len() == 0 {
				if current+uint(spacebuf.Len()) > width {
					current = 0
				} else {
					current += uint(spacebuf.Len())
					spacebuf.WriteTo(buf)
				}
				spacebuf.Reset()
			} else {
				current += uint(spacebuf.Len() + wordbuf.Len())
				spacebuf.WriteTo(buf)
				spacebuf.Reset()
				wordbuf.WriteTo(buf)
				wordbuf.Reset()
			}

			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) {
			if spacebuf.Len() == 0 || wordbuf.Len() > 0 {
				current += uint(spacebuf.Len() + wordbuf.Len())
				spacebuf.WriteTo(buf)
				spacebuf.Reset()
				wordbuf.WriteTo(buf)
				wordbuf.Reset()
			}

			spacebuf.WriteRune(char)
		} else {
			wordbuf.WriteRune(char)
			if current+uint(spacebuf.Len()+wordbuf.Len()) > width && uint(wordbuf.Len()) < width {
				buf.WriteString(br)
				current = 0
				spacebuf.Reset()
			}
		}
	}

	if wordbuf.Len() == 0 {
		if current+uint(spacebuf.Len()) <= width {
			spacebuf.WriteTo(buf)
		}
	} else {
		spacebuf.WriteTo(buf)
		wordbuf.WriteTo(buf)
	}

	return buf.String()
}

// StrLen — Get string length
func StrLen(str string) int {
	return len(str)
}

// MbStrLen — Get string length
func MbStrLen(str string) int {
	return utf8.RuneCountInString(str)
}

// StrRepeat — Repeat a string
func StrRepeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

// StrStr — Find the first occurrence of a string
func StrStr(haystack string, needle string) string {
	if needle == "" {
		return ""
	}

	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}

	return haystack[idx+len([]byte(needle))-1:]
}

// Strtr — Translate characters or replace substrings
// If the parameter length is 1, type is: map[string]string
// Strtr("baab", map[string]string{"ab": "01"}) will return "ba01"
// If the parameter length is 2, type is: string, string
// Strtr("baab", "ab", "01") will return "1001", a => 0; b => 1.
func Strtr(haystack string, params ...interface{}) string {
	ac := len(params)
	if ac == 1 {
		pairs := params[0].(map[string]string)
		length := len(pairs)
		if length == 0 {
			return haystack
		}

		oldnew := make([]string, length*2)
		for o, n := range pairs {
			if o == "" {
				return haystack
			}
			oldnew = append(oldnew, o, n)
		}

		return strings.NewReplacer(oldnew...).Replace(haystack)
	} else if ac == 2 {
		from := params[0].(string)
		to := params[1].(string)
		trlen, lt := len(from), len(to)

		if trlen > lt {
			trlen = lt
		}

		if trlen == 0 {
			return haystack
		}

		str := make([]uint8, len(haystack))
		var xlat [256]uint8
		var i int
		var j uint8

		if trlen == 1 {
			for i = 0; i < len(haystack); i++ {
				if haystack[i] == from[0] {
					str[i] = to[0]
				} else {
					str[i] = haystack[i]
				}
			}
			return string(str)
		}

		// trlen != 1
		for {
			xlat[j] = j
			if j++; j == 0 {
				break
			}
		}

		for i = 0; i < trlen; i++ {
			xlat[from[i]] = to[i]
		}

		for i = 0; i < len(haystack); i++ {
			str[i] = xlat[haystack[i]]
		}

		return string(str)
	}

	return haystack
}

// StrShuffle — Randomly shuffles a string
func StrShuffle(str string) string {
	runes := []rune(str)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := make([]rune, len(runes))

	for i, v := range r.Perm(len(runes)) {
		s[i] = runes[v]
	}

	return string(s)
}

// Trim — Strip whitespace (or other characters) from the beginning and end of a string
func Trim(str string, characterMask ...string) string {
	mask := ""

	if len(characterMask) == 0 {
		mask = " \\t\\n\\r\\0\\x0B"
	} else {
		mask = characterMask[0]
	}

	return strings.Trim(str, mask)
}

// Ltrim — Strip whitespace (or other characters) from the beginning of a string
func Ltrim(str string, characterMask ...string) string {
	mask := ""

	if len(characterMask) == 0 {
		mask = " \\t\\n\\r\\0\\x0B"
	} else {
		mask = characterMask[0]
	}

	return strings.TrimLeft(str, mask)
}

// Rtrim — Strip whitespace (or other characters) from the end of a string
func Rtrim(str string, characterMask ...string) string {
	mask := ""

	if len(characterMask) == 0 {
		mask = " \\t\\n\\r\\0\\x0B"
	} else {
		mask = characterMask[0]
	}

	return strings.TrimRight(str, mask)
}

// Explode — Split a string by a string
func Explode(delimiter, str string) []string {
	return strings.Split(str, delimiter)
}

// Implode — Join array elements with a string
func Implode(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)

	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}

	return buf.String()
}

// Chr — Generate a single-byte string from a number
func Chr(ascii int) string {
	return string(ascii)
}

// Ord — Convert the first byte of a string to a value between 0 and 255
func Ord(char string) int {
	r, _ := utf8.DecodeRune([]byte(char))
	return int(r)
}

// Nl2br — Inserts HTML line breaks before all newlines in a string
// \n\r, \r\n, \r, \n
func Nl2br(str string, isXhtml bool) string {
	r, n, runes := '\r', '\n', []rune(str)
	var br []byte
	if isXhtml {
		br = []byte("<br />")
	} else {
		br = []byte("<br>")
	}

	skip := false
	length := len(runes)
	var buf bytes.Buffer
	for i, v := range runes {
		if skip {
			skip = false
			continue
		}

		switch v {
		case n, r:
			if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
				buf.Write(br)
				skip = true
				continue
			}
			buf.Write(br)
		default:
			buf.WriteRune(v)
		}
	}

	return buf.String()
}

// JsonDecode — Decodes a JSON string
func JsonDecode(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}

// JsonEncode — Returns the JSON representation of a value
func JsonEncode(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}

// AddSlashes — Quote string with slashes
func AddSlashes(str string) string {
	var buf bytes.Buffer

	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}

		buf.WriteRune(char)
	}

	return buf.String()
}

// StripSlashes — Un-quotes a quoted string
func StripSlashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false

	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}

		buf.WriteRune(char)
	}

	return buf.String()
}

// QuoteMeta — Quote meta characters
func QuoteMeta(str string) string {
	var buf bytes.Buffer

	for _, char := range str {
		switch char {
		case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
			buf.WriteRune('\\')
		}

		buf.WriteRune(char)
	}

	return buf.String()
}

// HtmlEntities — Convert all applicable characters to HTML entities
func HtmlEntities(str string) string {
	return html.EscapeString(str)
}

// HtmlEntityDecode — Convert HTML entities to their corresponding characters
func HtmlEntityDecode(str string) string {
	return html.UnescapeString(str)
}

// Md5 — Calculate the md5 hash of a string
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// Md5File — Calculates the md5 hash of a given file
func Md5File(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Sha1 — Calculate the sha1 hash of a string
func Sha1(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha1File — Calculate the sha1 hash of a file
func Sha1File(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Crc32 — Calculates the crc32 polynomial of a string
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// Levenshtein — Calculate Levenshtein distance between two strings
// costIns: Defines the cost of insertion.
// costRep: Defines the cost of replacement.
// costDel: Defines the cost of deletion.
func Levenshtein(str1, str2 string, costIns, costRep, costDel int) int {
	var maxLen = 255
	l1 := len(str1)
	l2 := len(str2)

	if l1 == 0 {
		return l2 * costIns
	}

	if l2 == 0 {
		return l1 * costDel
	}

	if l1 > maxLen || l2 > maxLen {
		return -1
	}

	tmp := make([]int, l2+1)
	p1 := make([]int, l2+1)
	p2 := make([]int, l2+1)
	var c0, c1, c2 int
	var i1, i2 int

	for i2 := 0; i2 <= l2; i2++ {
		p1[i2] = i2 * costIns
	}

	for i1 = 0; i1 < l1; i1++ {
		p2[0] = p1[0] + costDel
		for i2 = 0; i2 < l2; i2++ {
			if str1[i1] == str2[i2] {
				c0 = p1[i2]
			} else {
				c0 = p1[i2] + costRep
			}

			c1 = p1[i2+1] + costDel
			if c1 < c0 {
				c0 = c1
			}

			c2 = p2[i2] + costIns
			if c2 < c0 {
				c0 = c2
			}

			p2[i2+1] = c0
		}

		tmp = p1
		p1 = p2
		p2 = tmp
	}

	c0 = p1[l2]

	return c0
}

// SimilarText — Calculate the similarity between two strings
func SimilarText(first, second string, percent *float64) int {
	var similarText func(string, string, int, int) int

	similarText = func(str1, str2 string, len1, len2 int) int {
		var sum, max int
		pos1, pos2 := 0, 0

		// Find the longest segment of the same section in two strings
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for l := 0; (i+l < len1) && (j+l < len2) && (str1[i+l] == str2[j+l]); l++ {
					if l+1 > max {
						max = l + 1
						pos1 = i
						pos2 = j
					}
				}
			}
		}

		if sum = max; sum > 0 {
			if pos1 > 0 && pos2 > 0 {
				sum += similarText(str1, str2, pos1, pos2)
			}

			if (pos1+max < len1) && (pos2+max < len2) {
				s1 := []byte(str1)
				s2 := []byte(str2)
				sum += similarText(string(s1[pos1+max:]), string(s2[pos2+max:]), len1-pos1-max, len2-pos2-max)
			}
		}

		return sum
	}

	l1, l2 := len(first), len(second)
	if l1+l2 == 0 {
		return 0
	}

	sim := similarText(first, second, l1, l2)
	if percent != nil {
		*percent = float64(sim*200) / float64(l1+l2)
	}

	return sim
}

// Soundex — Calculate the soundex key of a string
func Soundex(str string) string {
	if str == "" {
		panic("str: cannot be an empty string")
	}

	table := [26]rune{
		'0', '1', '2', '3', // A, B, C, D
		'0', '1', '2', // E, F, G
		'0',                          // H
		'0', '2', '2', '4', '5', '5', // I, J, K, L, M, N
		'0', '1', '2', '6', '2', '3', // O, P, Q, R, S, T
		'0', '1', // U, V
		'0', '2', // W, X
		'0', '2', // Y, Z
	}

	last, code, small := -1, 0, 0
	sd := make([]rune, 4)
	// build soundex string
	for i := 0; i < len(str) && small < 4; i++ {
		// ToUpper
		char := str[i]
		if char < '\u007F' && 'a' <= char && char <= 'z' {
			code = int(char - 'a' + 'A')
		} else {
			code = int(char)
		}

		if code >= 'A' && code <= 'Z' {
			if small == 0 {
				sd[small] = rune(code)
				small++
				last = int(table[code-'A'])
			} else {
				code = int(table[code-'A'])
				if code != last {
					if code != 0 {
						sd[small] = rune(code)
						small++
					}
					last = code
				}
			}
		}
	}

	// pad with "0"
	for ; small < 4; small++ {
		sd[small] = '0'
	}

	return string(sd)
}
