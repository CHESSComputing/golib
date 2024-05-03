package schema

import (
	"encoding/json"
	"fmt"
	"io"

	srvConfig "github.com/CHESSComputing/golib/config"
	services "github.com/CHESSComputing/golib/services"
)

var _httpReadRequest *services.HttpRequest
var Verbose int

// MetaDataDetails represents individual FOXDEN schema details dictionary
type MetaDataDetails struct {
	Schema       string            `json:"schema"`
	Units        map[string]string `json:"units"`
	Descriptions map[string]string `json:"descriptions"`
}

// MetaDataManager holds MetaDataDetails list for all FOXDEN schemas
type MetaDataManager struct {
	Records []MetaDataDetails
}

func (s *MetaDataManager) initManager() []MetaDataDetails {
	var records []MetaDataDetails
	if _httpReadRequest == nil {
		_httpReadRequest = services.NewHttpRequest("read", Verbose)
	}
	if s == nil {
		s = &MetaDataManager{}
		// fetch all schema details from upstream MetaData server
		rurl := fmt.Sprintf("%s/meta", srvConfig.Config.Services.MetaDataURL)
		if resp, err := _httpReadRequest.Get(rurl); err == nil {
			defer resp.Body.Close()
			if data, err := io.ReadAll(resp.Body); err == nil {
				if err := json.Unmarshal(data, &records); err == nil {
					s.Records = records
				}
			}
		}
	} else {
		records = s.Records
	}
	return records
}

// Units finds schema units map for given schema name
func (s *MetaDataManager) Units(sname string) map[string]string {
	records := s.initManager()
	for _, rec := range records {
		if rec.Schema != sname {
			continue
		}
		return rec.Units
	}
	empty := make(map[string]string)
	return empty
}

// Descriptions finds schema units map for given schema name
func (s *MetaDataManager) Descriptions(sname string) map[string]string {
	records := s.initManager()
	for _, rec := range records {
		if rec.Schema != sname {
			continue
		}
		return rec.Descriptions
	}
	empty := make(map[string]string)
	return empty
}
