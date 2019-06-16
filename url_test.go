package utils

import (
	"testing"
)

func TestUrl(t *testing.T) {
	tParseURL, _ := ParseUrl("http://username:password@hostname:9090/path?arg=value#anchor", -1)
	equal(t, map[string]string{"pass": "password", "path": "/path", "query": "arg=value", "fragment": "anchor", "scheme": "http", "host": "hostname", "port": "9090", "user": "username"}, tParseURL)

	tURLEncode := UrlEncode("http://golang.org?x y")
	equal(t, "http%3A%2F%2Fgolang.org%3Fx+y", tURLEncode)

	tURLDecode, _ := UrlDecode("http%3A%2F%2Fgolang.org%3Fx+y")
	equal(t, "http://golang.org?x y", tURLDecode)

	tRawurlencode := RawUrlEncode("http://golang.org?x y")
	equal(t, "http%3A%2F%2Fgolang.org%3Fx%20y", tRawurlencode)

	tRawurldecode, _ := RawUrlDecode("http%3A%2F%2Fgolang.org%3Fx%20y")
	equal(t, "http://golang.org?x y", tRawurldecode)

	tBase64Encode := Base64Encode("This is an encoded string")
	equal(t, "VGhpcyBpcyBhbiBlbmNvZGVkIHN0cmluZw==", tBase64Encode)

	tBase64Decode, _ := Base64Decode("VGhpcyBpcyBhbiBlbmNvZGVkIHN0cmluZw==")
	equal(t, "This is an encoded string", tBase64Decode)

	tHTTPBuildQuery := HttpBuildQuery(map[string][]string{"first": []string{"value"}, "multi": []string{"foo bar", "baz"}})
	equal(t, "first=value&multi=foo+bar&multi=baz", tHTTPBuildQuery)
}
