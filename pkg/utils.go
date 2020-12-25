package mmdb

import (
	"net"
	"strings"
)

func parseUint8(s string) uint8 {
	var n uint8

	for _, b := range []byte(s) {
		n *= uint8(10)
		n += uint8(b - '0')
	}

	return n
}

func parseUint32(s string) uint32 {
	var n uint32

	for _, b := range []byte(s) {
		n *= uint32(10)
		n += uint32(b - '0')
	}

	return n
}

func parseIPv4(ip string) net.IP {
	var octet, next string

	next = ip
	octets := make(net.IP, net.IPv4len)
	for i := 0; i < net.IPv4len; i++ {
		n := strings.IndexByte(next, '.')
		if n == -1 {
			octets[i] = byte(parseUint8(next))
			break
		}

		octet, next = next[:n], next[n+1:]
		octets[i] = byte(parseUint8(octet))
	}

	return octets
}
