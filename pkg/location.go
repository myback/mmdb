package mmdb

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type ids map[uint32]string

func (idx *ids) put(id, value string) {
	(*idx)[parseUint32(id)] = value
}

func (idx *ids) Get(id string) (string, bool) {
	v, ok := (*idx)[parseUint32(id)]
	return v, ok
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

		// TODO: remove kludge
		if v == "" {
			v = record[2]
		}

		if k == "" || v == "" {
			fmt.Printf("[WARN] %s empty data; k:(%d:%s); v:(%d:%s); %#v\n", filename, keyIdx, k, valueIdx, v, record)
		} else {
			out.put(k, v)
		}
	}

	return nil
}
