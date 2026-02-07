package services

import (
	"encoding/json"
	"fmt"
)

// Record define Mongo record
// type Record map[string]interface{}

// MetaRecord represents meta-data record used for injection
type MetaRecord struct {
	Schema string
	Record map[string]any
}

// String converts ServiceResponse into string representation
func (s *MetaRecord) JsonString() string {
	data, _ := json.MarshalIndent(s, "", "  ")
	return string(data)
}

// ServiceQuery represents service query along with its results
type ServiceQuery struct {
	Query      string         `json:"query"`
	Spec       map[string]any `json:"spec"`
	Projection map[string]any `json:"projection"`
	SQL        string         `json:"sql"`
	Idx        int            `json:"idx"`
	Limit      int            `json:"limit"`
	SortKeys   []string       `json:"sort_keys"`
	SortOrder  int            `json:"sort_order"`
}

// ServiceResults represents service results
type ServiceResults struct {
	NRecords int              `json:"nrecords"`
	Records  []map[string]any `json:"records"`
}

// ServiceRequest represents service request structure
type ServiceRequest struct {
	Client       string       `json:"client"`
	ServiceQuery ServiceQuery `json:"service_query"`
}

// String converts ServiceRequest into string representation
func (s *ServiceRequest) String() string {
	out := fmt.Sprintf("\n\tclient    : %s", s.Client)
	out += fmt.Sprintf("\n\tquery     : %s", s.ServiceQuery.Query)
	out += fmt.Sprintf("\n\tspec      : %+v", s.ServiceQuery.Spec)
	out += fmt.Sprintf("\n\tprojection: %+v", s.ServiceQuery.Projection)
	out += fmt.Sprintf("\n\tsql       : %s", s.ServiceQuery.SQL)
	out += fmt.Sprintf("\n\tidx       : %d", s.ServiceQuery.Idx)
	out += fmt.Sprintf("\n\tlimit     : %d", s.ServiceQuery.Limit)
	out += fmt.Sprintf("\n\tsort key  : %v", s.ServiceQuery.SortKeys)
	out += fmt.Sprintf("\n\tsort order: %d", s.ServiceQuery.SortOrder)
	return out
}

// ServiceResponse represents service response structure
type ServiceResponse struct {
	HttpCode     int            `json:"http_code"`
	SrvCode      int            `json:"service_code"`
	Service      string         `json:"service"`
	Status       string         `json:"status"`
	Error        string         `json:"error"`
	ServiceQuery ServiceQuery   `json:"service_query,omitempty"`
	Results      ServiceResults `json:"results,omitempty"`
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

// HtmlString converts ServiceResponse into string representation for html
func (s *ServiceResponse) HtmlString() string {
	var out string
	out += fmt.Sprintf("<b>Service:</b> %s<br/>", s.Service)
	out += fmt.Sprintf("<b>Code:</b> %d<br/>", s.SrvCode)
	out += fmt.Sprintf("<b>Status:</b> %s<br/>", s.Status)
	out += fmt.Sprintf("<b>Error:</b> %s<br/>", s.Error)
	out += fmt.Sprintf("<b>Timestamp:</b> %s<br/>", s.Timestamp)
	return out
}

// JsonString converts ServiceResponse into string representation
func (s *ServiceResponse) JsonString() string {
	data, _ := json.MarshalIndent(s, "", "  ")
	return string(data)
}

// JsonBytes converts ServiceResponse into bytes representation
func (s *ServiceResponse) JsonBytes() []byte {
	data, _ := json.MarshalIndent(s, "", "  ")
	return data
}
