package ql

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// ServiceMap defines FOXDEN service QL mapping
type ServiceMap map[string][]string

type QLManager struct {
	Map ServiceMap
}

// Load function loads service map from given file name
func (q *QLManager) Load(fname string) error {
	if q.Map == nil {
		q.Map = make(ServiceMap)
	}
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var srvMap ServiceMap
	err = json.Unmarshal(data, &srvMap)
	if err != nil {
		return err
	}
	q.Map = srvMap
	return nil
}

// Keys provides list of keys associated with FOXDEN service name
func (q *QLManager) Keys(srv string) []string {
	var keys []string
	if val, ok := q.Map[srv]; ok {
		return val
	}
	return keys
}

// Services returns list of services known to QL manager
func (q *QLManager) Services() []string {
	var srv []string
	for k, _ := range q.Map {
		srv = append(srv, k)
	}
	return srv
}
