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
	Map ServiceMap
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
