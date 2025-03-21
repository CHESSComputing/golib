package doi

import (
	srvConfig "github.com/CHESSComputing/golib/config"
	datacite "github.com/CHESSComputing/golib/datacite"
)

// DataciteProvider represents Datacite provider
type DataciteProvider struct {
	Name    string
	Verbose int
}

// Init function initializes Datacite publisher
func (d *DataciteProvider) Init() {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	d.Name = "foxden-datacite"
}

// Publish provides publication of dataset with did and description
func (d *DataciteProvider) Publish(did, description string, record map[string]any, publish bool) (string, string, error) {
	doi, doiLink, err := datacite.Publish(did, description, record, publish)
	return doi, doiLink, err
}
