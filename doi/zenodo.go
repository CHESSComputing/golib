package doi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/services"
	"github.com/CHESSComputing/golib/zenodo"
)

// ZenodoProvider represents Zenodo provider
type ZenodoProvider struct {
	HttpRequest *services.HttpRequest
	Token       string
}

// Init function initializes Zenodo provider
func (z *ZenodoProvider) Init() {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	if z.HttpRequest == nil {
		z.HttpRequest = services.NewHttpRequest("read", 0)
		z.HttpRequest.Token = z.Token
	}
}

// Publish provides publication of dataset with did and description
func (z *ZenodoProvider) Publish(did, description string) (string, string, error) {
	var doi, doiLink string
	var err error
	docId, err := zenodo.CreateRecord()
	if err != nil {
		return doi, doiLink, err
	}

	// extract meta-data record for our did
	query := fmt.Sprintf("{\"did\": \"%s\"}", did)
	rec := services.ServiceRequest{
		Client:       "foxden-doi",
		ServiceQuery: services.ServiceQuery{Query: query, Idx: 0, Limit: -1},
	}

	data, err := json.Marshal(rec)
	rurl := fmt.Sprintf("%s/search", srvConfig.Config.Services.MetaDataURL)
	resp, err := z.HttpRequest.Post(rurl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return doi, doiLink, err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return doi, doiLink, err
	}
	var records []map[string]any
	err = json.Unmarshal(data, &records)
	if err != nil {
		return doi, doiLink, err
	}

	// add foxden record
	frec := zenodo.FoxdenRecord{Did: did, MetaData: records}
	err = zenodo.AddRecord(docId, "foxden-metadata.json", frec)
	if err != nil {
		return doi, doiLink, err
	}

	// create new meta-data record
	creator := zenodo.Creator{Name: "FOXDEN", Affiliation: "Cornell University"}
	mrec := zenodo.MetaDataRecord{
		PublicationType: "deliverable",
		UploadType:      "dataset",
		Description:     description,
		Keywords:        []string{"FOXDEN"},
		Title:           fmt.Sprintf("FOXDEN dataset did=%s", did),
		Licences:        []string{"MIT"},
		Creators:        []zenodo.Creator{creator},
	}

	err = zenodo.UpdateRecord(docId, mrec)
	if err != nil {
		return doi, doiLink, err
	}

	// publish record
	doiRecord, err := zenodo.PublishRecord(docId)
	if err != nil {
		return doi, doiLink, err
	}
	return doiRecord.Doi, doiRecord.DoiUrl, nil
}
