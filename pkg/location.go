package mmdb

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type ids map[uint32]string

func parseUint32(s string) (uint32, error) {
	_id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return uint32(0), err
	}

	return uint32(_id), nil
}

func (idx *ids) put(id, value string) error {
	_id, err := parseUint32(id)
	if err == nil {
		(*idx)[_id] = value
	}

	return err
}

func (idx *ids) Get(id string) (string, bool, error) {
	_id, err := parseUint32(id)
	if err != nil {
		return "", false, err
	}

	v, ok := (*idx)[_id]
	return v, ok, nil
}

// NewIDs load IDs from CSV file
func NewIDs(filename, csvKeyName, csvValueName string) (*ids, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	parser := csv.NewReader(file)
	db := ids{}

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
			db.put(k, v)
		}
	}

	return &db, nil
}
