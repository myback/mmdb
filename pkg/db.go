package mmdb

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
)

type ipdb map[uint32]string

type Getter interface {
	Get(string) string
}

func ip2uint32(ip net.IP) (u uint32) {
	binary.Read(bytes.NewBuffer(ip), binary.BigEndian, &u)

	return u
}

func (db *ipdb) put(ip, value string) {
	k := ip2uint32(net.ParseIP(ip).To4())
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

		if v, ok := (*db)[ip2uint32(netIP)]; ok {
			return v
		}
	}

	return ""
}

// NewDB return ipdb object
func NewDB(filename, csvKeyName, csvValueName string, ids *ids) (*ipdb, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	parser := csv.NewReader(file)
	db := ipdb{}

	firstLine := true
	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var keyIdx, valueIdx int
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

			if keyIdx == 0 || valueIdx == 0 {
				return nil, fmt.Errorf("key_name(%d) or value_name(%d) invalid", keyIdx, valueIdx)
			}

			continue
		}

		k := record[keyIdx]
		v := record[valueIdx]
		if v == "" {
			v = record[valueIdx-1]
		}

		if k == "" || v == "" {
			fmt.Printf("[WARN] empty data; k:%s; v:%s; %#v", k, v, record)
		} else {
			v, ok, err := ids.Get(v)
			if err != nil {
				fmt.Printf("[WARN] err: %s", err)
			}
			if ok {
				db.put(k, v)
			} else {
				fmt.Printf("[WARN] ids not found for id: %s", v)
			}
		}
	}

	return &db, nil
}
