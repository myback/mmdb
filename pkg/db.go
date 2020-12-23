package mmdb

import (
	"bytes"
	"encoding/binary"
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

func IP2uint32(ip net.IP) (u uint32) {
	binary.Read(bytes.NewBuffer(ip), binary.BigEndian, &u)

	return u
}

func (db *ipdb) put(ip, value string) {
	ip = strings.Split(ip, "/")[0]
	k := IP2uint32(net.ParseIP(ip).To4())
	(*db)[k] = value
}

func (db *ipdb) Get(ip string) string {
	var prevIP net.IP
	for i := 32; i > 1; i-- {
		_, ipnet, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip, i))
		netIP := ipnet.IP.To4()

		if netIP.Equal(prevIP) {
			continue
		}

		prevIP = netIP

		if v, ok := (*db)[IP2uint32(netIP)]; ok {
			return v
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

		if k == "" || v == "" {
			out.put(k, "unknown")
			fmt.Printf("[WARN] UNKNOWN: %s\n", record[0])
		} else {
			val, ok, err := ids.Get(v)
			if err != nil {
				fmt.Printf("[WARN] err: %d %s\n", valueIdx, err)
			}

			if ok {
				out.put(k, val)
			} else {
				fmt.Printf("[WARN] ids %s not found for id: %s\n", filename, v)
			}
		}
	}

	return nil
}
