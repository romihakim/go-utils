package utils

import (
	"encoding/binary"
	"net"
	"os"
	"strings"
)

// GetHostName — Gets the host name
func GetHostName() (string, error) {
	return os.Hostname()
}

// GetHostByName — Get the IPv4 address corresponding to a given Internet host name
func GetHostByName(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		for _, v := range ips {
			if v.To4() != nil {
				return v.String(), nil
			}
		}

		return "", nil
	}

	return "", err
}

// GetHostByNamel — Get a list of IPv4 addresses corresponding to a given Internet host name
func GetHostByNamel(hostname string) ([]string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		var ipstrs []string
		for _, v := range ips {
			if v.To4() != nil {
				ipstrs = append(ipstrs, v.String())
			}
		}
		return ipstrs, nil
	}
	return nil, err
}

// GetHostByAddr — Get the Internet host name corresponding to a given IP address
func GetHostByAddr(ipAddress string) (string, error) {
	names, err := net.LookupAddr(ipAddress)
	if names != nil {
		return strings.TrimRight(names[0], "."), nil
	}

	return "", err
}

// Ip2Long — Converts a string containing an (IPv4) Internet Protocol dotted address into a long integer
// IPv4
func Ip2Long(ipAddress string) uint32 {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}

	return binary.BigEndian.Uint32(ip.To4())
}

// Long2Ip — Converts an long integer address into a string in (IPv4) Internet standard dotted format
// IPv4
func Long2Ip(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	ip := net.IP(ipByte)

	return ip.String()
}
