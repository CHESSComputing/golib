package ql

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// QLRecord defines structure of QL key
type QLRecord struct {
	Key         string `json:"key"`
	Description string `json:"description,omitempty"`
	Service     string `json:"service"`
	Units       string `json:"units,omitempty"`
	Schema      string `json:"schema,omitempty"`
	DataType    string `json:"type"`
}

// FillEmpty assign N/A to empty attributes
func (q *QLRecord) FillEmpty() {
	if q.Description == "" {
		q.Description = "N/A"
	}
	if q.Schema == "" {
		q.Schema = "N/A"
	}
	if q.Units == "" {
		q.Units = "N/A"
	}
	if q.DataType == "" {
		q.DataType = "N/A"
	}
}

// String returns string representation of ql key
func (q *QLRecord) String() string {
	q.FillEmpty()
	out := fmt.Sprintf("Key:%s Description:%s Service:%s Schema:%s Units:%s DataType:%s",
		q.Key, q.Description, q.Service, q.Schema, q.Units, q.DataType)
	return out
}

// Details function provides record representation
func (q *QLRecord) Details(show string) string {
	q.FillEmpty()
	repr := fmt.Sprintf("%s, service: %s, schema: %s, units: %s, data-type: %s",
		q.Description, q.Service, q.Schema, q.Units, q.DataType)
	if show == "key" {
		return q.Key
	}
	if show == "description" {
		return q.Description
	}
	if show == "service" {
		return q.Service
	}
	if show == "schema" {
		return q.Schema
	}
	if show == "units" {
		return q.Units
	}
	if show == "data-type" {
		return q.DataType
	}
	return repr
}

// QLRecords return list of ql keys
func QLRecords(qlKey string) ([]QLRecord, error) {
	var records []QLRecord
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	fname := srvConfig.Config.QL.ServiceMapFile
	file, err := os.Open(fname)
	if err != nil {
		return records, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return records, err
	}
	err = json.Unmarshal(data, &records)
	return records, err
}
