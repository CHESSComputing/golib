package server

import (
	"github.com/CHESSComputing/golib/mongo"
)

// MetaRecord represents meta-data record used for injection
type MetaRecord struct {
	Schema string
	Record mongo.Record
}

// MetaResponse represents response struct from meta-data service
type MetaResponse struct {
	Query    string
	Spec     any
	Idx      int
	Limit    int
	NRecords int
	Records  []mongo.Record
}

// DBSResponse represents reponse struct from dbs service
type DBSResponse struct {
	Query    string
	SQL      string
	Idx      int
	Limit    int
	NRecords int
	Records  []mongo.Record
}

// DiscoveryResponse represents response struct from discovery service
type DiscoveryResponse struct {
	DBSResponse  `json:"dbs"`
	MetaResponse `json:"meta"`
}

// StatusStatus represents status record
type ServiceStatus struct {
	HttpCode int
	Code     int
	Service  string
	Status   string
	Error    error
}
