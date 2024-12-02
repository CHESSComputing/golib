package ql

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	utils "github.com/CHESSComputing/golib/utils"
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
func (q *QLManager) ServiceQueries(query string) (map[string]map[string]any, error) {
	sqMap := make(map[string]map[string]any)
	spec, err := ParseQuery(query)
	if err != nil {
		return nil, err
	}
	for key, smap := range spec {
		for srv, _ := range q.Map {
			if allowed := q.QueryKeyAllowed(key, srv); allowed {
				if val, ok := sqMap[srv]; ok {
					val[key] = smap
					sqMap[srv] = val
				} else {
					sqMap[srv] = map[string]any{key: smap}
				}
			}
		}
	}
	return sqMap, nil
}

// Determines if a key from a user query is allowed for querying the given service.
// A user query key is considered "allowed" if it is an exact match for one of the
// service's allowed query keys OR if the user query key has a prefix equal to one
// of the allowed service query keys followed by a ".". The latter matching condition
// allows for queries on nested fields with dot notation.
func (q *QLManager) QueryKeyAllowed(key string, service string) bool {
	if service_keys, ok := q.Map[service]; ok {
		if utils.InList(key, service_keys) {
			return true
		}
		for _, service_key := range service_keys {
			if strings.HasPrefix(key, fmt.Sprintf("%s.", service_key)) {
				return true
			}
		}
	}
	return false
}
