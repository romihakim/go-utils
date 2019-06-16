package utils

import (
	"testing"
)

func TestNetwork(t *testing.T) {
	tGethostname, _ := GetHostName()
	gt(t, float64(len(tGethostname)), 0)

	tIP2long := Ip2Long("8.8.8.8")
	equal(t, uint32(134744072), tIP2long)

	tLong2ip := Long2Ip(134744072)
	equal(t, "8.8.8.8", tLong2ip)

	tGethostbyname, _ := GetHostByName("localhost")
	equal(t, "127.0.0.1", tGethostbyname)

	tGethostbynamel, _ := GetHostByNamel("localhost")
	gt(t, float64(len(tGethostbynamel)), 0)

	tGethostbyaddr, _ := GetHostByAddr("127.0.0.1")
	equal(t, "localhost", tGethostbyaddr)
}
