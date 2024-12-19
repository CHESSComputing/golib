package doi

import (
	"fmt"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/zenodo"
)

// ZenodoProvider represents Zenodo provider
type ZenodoProvider struct {
}

// Init function initializes Zenodo provider
func (z *ZenodoProvider) Init() {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
}

// Publish provides publication of dataset with did and description
func (z *ZenodoProvider) Publish(did, description string, record any) (string, string, error) {
	var doi, doiLink string
	var err error
	docId, err := zenodo.CreateRecord()
	if err != nil {
		return doi, doiLink, err
	}

	// add foxden record
	frec := zenodo.FoxdenRecord{Did: did, MetaData: record}
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
