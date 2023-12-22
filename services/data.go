package server

import (
	"encoding/json"
	"fmt"

	"github.com/CHESSComputing/golib/mongo"
)

// MetaRecord represents meta-data record used for injection
type MetaRecord struct {
	Schema string
	Record mongo.Record
}

// String converts ServiceResponse into string representation
func (s *MetaRecord) JsonString() string {
	data, _ := json.MarshalIndent(s, "", "  ")
	return string(data)
}

// ServiceQuery represents service query along with its results
type ServiceQuery struct {
	Query string `json:"query"`
	Spec  any    `json:"spec"`
	SQL   string `json:"sql"`
	Idx   int    `json:"idx"`
	Limit int    `json:"limit"`
}

// ServiceResults represents service results
type ServiceResults struct {
	NRecords int            `json:"nrecords"`
	Records  []mongo.Record `json:"records"`
}

// ServiceRequest represents service request structure
type ServiceRequest struct {
	Client       string       `json:"client"`
	ServiceQuery ServiceQuery `json:"service_query"`
}

// String converts ServiceRequest into string representation
func (s *ServiceRequest) String() string {
	data, _ := json.MarshalIndent(s, "", "  ")
	return string(data)
}

// ServiceResponse represents service response structure
type ServiceResponse struct {
	HttpCode     int            `json:"http_code"`
	SrvCode      int            `json:"service_code"`
	Service      string         `json:"service"`
	Status       string         `json:"status"`
	Error        string         `json:"error"`
	ServiceQuery ServiceQuery   `json:"service_query"`
	Results      ServiceResults `json:"results"`
	Timestamp    string         `json:"timestamp"`
}

// String converts ServiceResponse into string representation
func (s *ServiceResponse) String() string {
	var out string
	out += fmt.Sprintf("Service     : %s\n", s.Service)
	out += fmt.Sprintf("Code        : %d\n", s.SrvCode)
	out += fmt.Sprintf("Status      : %s\n", s.Status)
	out += fmt.Sprintf("Error       : %s\n", s.Error)
	out += fmt.Sprintf("Timestamp   : %s\n", s.Timestamp)
	return out
}

// String converts ServiceResponse into string representation
func (s *ServiceResponse) JsonString() string {
	data, _ := json.MarshalIndent(s, "", "  ")
	return string(data)
}
