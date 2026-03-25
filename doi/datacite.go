package doi

import (
	"fmt"

	datacite "github.com/CHESSComputing/golib/datacite"
)

// DataciteProvider represents Datacite provider
type DataciteProvider struct {
	Name    string
	Verbose int
}

// Init function initializes Datacite publisher
func (d *DataciteProvider) Init() {
	d.Name = "foxden-datacite"
}

// Publish provides publication of dataset with did and description
func (d *DataciteProvider) Publish(authors []string, did, description string, record map[string]any, publish bool) (string, string, error) {
	doi, doiLink, err := datacite.Publish(authors, did, description, record, publish, d.Verbose)
	if err != nil {
		return doi, doiLink, fmt.Errorf("[golib.doi.DataciteProvider.Publish] datacite.Publish error: %w", err)
	}
	return doi, doiLink, nil
}

// MakePublic provides publication of draft DOI
func (d *DataciteProvider) MakePublic(doi string) error {
	return datacite.MakePublic(doi, d.Verbose)
}
