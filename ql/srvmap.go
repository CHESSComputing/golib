package ql

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sort"

	utils "github.com/CHESSComputing/golib/utils"
	bson "go.mongodb.org/mongo-driver/bson"
)

// ServiceMap defines FOXDEN service QL mapping
type ServiceMap map[string][]string

// QLManager represents QL manager
type QLManager struct {
	Map     ServiceMap
	Records []QLRecord
}

// Init function loads service map from given file name
func (q *QLManager) Init(fname string) error {
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

	// each record in QL.ServiceMapFile has the following form:[QLRecord1, QLRecord2]
	srvMap := make(ServiceMap)
	var records []QLRecord
	err = json.Unmarshal(data, &records)
	if err != nil {
		return err
	}
	for _, rec := range records {
		// collect ql kys for each service
		if qlKeys, ok := srvMap[rec.Service]; ok {
			qlKeys = append(qlKeys, rec.Key)
			srvMap[rec.Service] = qlKeys
		} else {
			srvMap[rec.Service] = []string{rec.Key}
		}
	}
	q.Map = srvMap
	q.Records = records
	return nil
}

// Keys provides list of keys associated with FOXDEN service name
func (q *QLManager) Keys(srv string) []string {
	var keys []string
	if val, ok := q.Map[srv]; ok {
		sort.Strings(val)
		return val
	}
	sort.Strings(keys)
	return keys
}

// Services returns list of services known to QL manager
func (q *QLManager) Services() []string {
	var srv []string
	for k, _ := range q.Map {
		srv = append(srv, k)
	}
	sort.Strings(srv)
	return srv
}

// ServiceQueries parses given query string into list of service queries
func (q *QLManager) ServiceQueries(query string) (map[string]bson.M, error) {
	sqMap := make(map[string]bson.M)
	spec, err := ParseQuery(query)
	if err != nil {
		return nil, err
	}
	for key, smap := range spec {
		for srv, skeys := range q.Map {
			if utils.InList(key, skeys) {
				if val, ok := sqMap[srv]; ok {
					val[key] = smap
					sqMap[srv] = val
				} else {
					sqMap[srv] = bson.M{key: smap}
				}
			}
		}
	}
	return sqMap, nil
}
