package ql

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// ServiceMap defines FOXDEN service QL mapping
type ServiceMap map[string][]string

// Load function loads service map from given file name
func (s *ServiceMap) Load(fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	var srvMap ServiceMap
	err = json.Unmarshal(data, &srvMap)
	if err != nil {
		return nil
	}
	s = &srvMap
	return nil
}

// Keys provides list of keys associated with FOXDEN service name
func (s *ServiceMap) Keys(srv string) []string {
	var keys []string
	smap := *s
	if val, ok := smap[srv]; ok {
		return val
	}
	return keys
}
