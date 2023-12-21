package server

import (
	"github.com/CHESSComputing/golib/mongo"
)

// MetaRecord represents meta-data record used for injection
type MetaRecord struct {
	Schema string
	Record mongo.Record
}

// ServiceData represents service query along with its results
type ServiceData struct {
	Query    string
	Spec     any
	SQL      string
	Idx      int
	Limit    int
	NRecords int
	Records  []mongo.Record
}

// StatusStatus represents status record
type ServiceResponse struct {
	HttpCode  int `json:"http_code"`
	SrvCode   int `json:"service_code"`
	Service   string
	Status    string
	Error     error
	Data      ServiceData
	Timestamp string
}
