package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// MetaDataRecords return meta-data records for given query
func MetaDataRecords(query string, skeys []string, sorder, idx, limit int) ([]map[string]any, error) {
	var records []map[string]any
	rec := ServiceRequest{
		Client:       "foxden",
		ServiceQuery: ServiceQuery{Query: query, Idx: idx, Limit: limit, SortKeys: skeys, SortOrder: sorder},
	}

	data, err := json.Marshal(rec)
	if err != nil {
		return records, fmt.Errorf("[golib.services.MetaDataRecords] json.Marshal error: %w", err)
	}
	rurl := fmt.Sprintf("%s/search", srvConfig.Config.Services.MetaDataURL)
	httpReadRequest := NewHttpRequest("read", 0)
	httpReadRequest.GetToken()
	resp, err := httpReadRequest.Post(rurl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return records, fmt.Errorf("[golib.services.MetaDataRecords] httpReadRequest.Post error: %w", err)
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return records, fmt.Errorf("[golib.services.MetaDataRecords] io.ReadAll error: %w", err)
	}
	err = json.Unmarshal(data, &records)
	if err != nil {
		return records, nil
	}
	return records, nil
}
