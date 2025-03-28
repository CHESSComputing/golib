package doi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	datacite "github.com/CHESSComputing/golib/datacite"
	"github.com/CHESSComputing/golib/zenodo"
)

// ZenodoProvider represents Zenodo provider
type ZenodoProvider struct {
	Verbose int
}

// Init function initializes Zenodo provider
func (z *ZenodoProvider) Init() {
}

// Publish provides publication of dataset with did and description
func (z *ZenodoProvider) Publish(did, description string, record map[string]any, publish bool) (string, string, error) {
	var doi, doiLink string
	var err error

	// create new meta-data record
	creator := zenodo.Creator{Name: "FOXDEN", Affiliation: "Cornell University"}
	rec := zenodo.MetaDataRecord{
		PublicationType: "deliverable",
		UploadType:      "dataset",
		Description:     description,
		Keywords:        []string{"FOXDEN"},
		Title:           fmt.Sprintf("FOXDEN did=%s", did),
		Licences:        []string{"MIT"},
		Creators:        []zenodo.Creator{creator},
		PreserveDoi:     true,
	}
	mrec := make(map[string]any)
	mrec["metadata"] = rec
	payload, err := json.Marshal(mrec)
	if err != nil {
		return doi, doiLink, err
	}

	// create new Zenodo record
	doc, err := zenodo.CreateRecord(payload)
	if err != nil {
		return doi, doiLink, err
	}
	docId := doc.Id
	if docId == 0 {
		log.Println("ERROR: unable to create Zenodo document, docId=0")
		return doi, doiLink, errors.New("unable to create Zenodo document, docId=0")
	}
	doi = doc.MetaData.PrereserveDoi.Doi
	if doi != "" {
		doiLink = fmt.Sprintf("https://doi.org/%s", doi)
	}
	if z.Verbose > 0 {
		log.Printf("Created new Zenodo record docId=%v doi=%v", docId, doi)
	}

	// add foxden datacite metadata record
	frec := zenodo.FoxdenRecord{Did: did, MetaData: record}
	if payload, err := datacite.DataCiteMetadata(did, description, record, publish); err == nil {
		var rec map[string]any
		if err := json.Unmarshal(payload, &rec); err == nil {
			frec = zenodo.FoxdenRecord{Did: did, MetaData: rec}
		}
	}
	err = zenodo.AddRecord(docId, "foxden-datacite.json", frec)
	if err != nil {
		log.Println("ERROR: unable to add foxden-datacite.json record", err)
		return doi, doiLink, err
	}

	if !publish {
		log.Println("Zenodo record has been created with docId, but it is not published", docId)
		return doi, doiLink, nil
	}

	// publish record
	doiRecord, err := zenodo.PublishRecord(docId)
	if err != nil {
		return doi, doiLink, err
	}
	if z.Verbose > 0 {
		log.Printf("Published doi record %+v", doiRecord)
	}
	return doiRecord.Doi, doiRecord.DoiUrl, nil
}

// MakePublic provides publication of draft DOI
func (m *ZenodoProvider) MakePublic(doi string) error {
	return zenodo.MakePublic(doi)
}
