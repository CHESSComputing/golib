package server

import (
	"github.com/CHESSComputing/golib/mongo"
)

// MetaRecord represents meta-data record used for injection
type MetaRecord struct {
	Schema string
	Record mongo.Record
}

// ServiceQuery represents service query along with its results
type ServiceQuery struct {
	Query string
	Spec  any
	SQL   string
	Idx   int
	Limit int
}

// ServiceResults represents service results
type ServiceResults struct {
	NRecords int
	Records  []mongo.Record
}

// ServiceRequest represents service request structure
type ServiceRequest struct {
	Client string
	User   string
	Query  ServiceQuery
}

// ServiceResponse represents service response structure
type ServiceResponse struct {
	HttpCode  int `json:"http_code"`
	SrvCode   int `json:"service_code"`
	Service   string
	Status    string
	Error     error
	Query     ServiceQuery
	Results   ServiceResults
	Timestamp string
}
