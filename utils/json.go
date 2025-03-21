package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// ReadJson reads json data from a given file name and returns formatted JSON in bytes
func ReadJson(fname string) []byte {
	if _, err := os.Stat(fname); err == nil {
		file, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		if bdata, err := io.ReadAll(file); err == nil {
			var data map[string]any
			if err := json.Unmarshal(bdata, &data); err == nil {
				if fdata, err := json.MarshalIndent(data, "", "  "); err == nil {
					return fdata
				}

			}
		}
	}
	return []byte{}
}
