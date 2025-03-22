package doi

import (
	"fmt"
	"log"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/zenodo"
)

// ZenodoProvider represents Zenodo provider
type ZenodoProvider struct {
	Verbose int
}

// Init function initializes Zenodo provider
func (z *ZenodoProvider) Init() {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
}

// Publish provides publication of dataset with did and description
func (z *ZenodoProvider) Publish(did, description string, record map[string]any, publish bool) (string, string, error) {
	var doi, doiLink string
	var err error
	docId, err := zenodo.CreateRecord()
	if err != nil {
		return doi, doiLink, err
	}
	if z.Verbose > 0 {
		log.Println("Created zenodo record", docId)
	}

	// add foxden record
	frec := zenodo.FoxdenRecord{Did: did, MetaData: record}
	err = zenodo.AddRecord(docId, "foxden-metadata.json", frec)
	if err != nil {
		return doi, doiLink, err
	}
	if z.Verbose > 0 {
		log.Println("Created foxden record")
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
	if z.Verbose > 0 {
		log.Println("Updated doi record")
	}

	// publish record
	doiRecord, err := zenodo.PublishRecord(docId)
	if err != nil {
		return doi, doiLink, err
	}
	if z.Verbose > 0 {
		log.Println("Published doi record")
	}
	return doiRecord.Doi, doiRecord.DoiUrl, nil
}
