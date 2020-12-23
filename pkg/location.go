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
func NewIDs(filename, csvKeyName, csvValueName, csvValueBak string, out *ids) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	firstLine := true
	var keyIdx, valueIdx, valueIdxBak int

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
				case csvValueBak:
					valueIdxBak = i
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
			v = record[valueIdxBak]
		}

		if v == "" {
			v = record[2]
		}

		if k == "" || v == "" {
			fmt.Printf("[WARN] %s empty data; k:(%d:%s); v:(%d:%s); %#v\n", filename, keyIdx, k, valueIdx, v, record)
		} else {
			if err := out.put(k, v); err != nil {
				fmt.Println("[ERR]", err)
			}
		}
	}

	return nil
}
