package doi

import (
	"fmt"

	"github.com/CHESSComputing/golib/zenodo"
)

type Zenodo struct {
}

func (z *Zenodo) Publish(did, description string) (string, string, error) {
	var doi, doiLink string
	var err error
	docId, err := zenodo.CreateRecord()
	if err != nil {
		return doi, doiLink, err
	}

	// add foxden record
	// TODO: add to FoxdenRecord MetaData
	frec := zenodo.FoxdenRecord{Beamline: "test-beamline", Type: "raw-data", MetaData: "todo"}
	err = zenodo.AddRecord(docId, "foxden-meta.json", frec)

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
