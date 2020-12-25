package mmdb

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type ipdb map[uint32]string

type Getter interface {
	Get(string) string
}

func IP2uint32(ip net.IP) uint32 {
	return uint32(ip[3]) | uint32(ip[2])<<8 | uint32(ip[1])<<16 | uint32(ip[0])<<24
}

func (db *ipdb) put(ip, value string) {
	i := strings.IndexByte(ip, '/')
	k := IP2uint32(parseIPv4(ip[:i]))
	(*db)[k] = value
}

func (db *ipdb) Get(s string) string {
	var prevIP string

	ip := parseIPv4(s)
	for i := 32; i > 1; i-- {
		netIP := ip.Mask(net.CIDRMask(i, 8*net.IPv4len))

		netIPstr := netIP.String()
		if prevIP != netIPstr {
			if v, ok := (*db)[IP2uint32(netIP)]; ok {
				return v
			}

			prevIP = netIPstr
		}
	}

	return ""
}

// NewDB return ipdb object
func NewDB(filename, csvKeyName, csvValueName string, ids *ids, out *ipdb) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	firstLine := true
	var keyIdx, valueIdx int

	parser := csv.NewReader(file)
	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if firstLine {
			firstLine = false

			for i, v := range record {
				switch v {
				case csvKeyName:
					keyIdx = i
				case csvValueName:
					valueIdx = i
				}
			}

			if record[keyIdx] != csvKeyName || record[valueIdx] != csvValueName {
				return fmt.Errorf("key_name(%d:%s) or value_name(%d:%s) invalid", keyIdx, csvKeyName, valueIdx, csvValueName)
			}

			continue
		}

		k := record[keyIdx]
		v := record[valueIdx]
		if v == "" {
			v = record[valueIdx+1]
		}

		val, ok := ids.Get(v)
		if !ok {
			val = "unknown"
			fmt.Printf("[WARN] UNKNOWN: %s\n", record[0])
		}

		out.put(k, val)

	}

	return nil
}
