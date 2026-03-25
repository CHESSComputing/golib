package doi

import (
	"fmt"

	srvConfig "github.com/CHESSComputing/golib/config"
	materialscommons "github.com/CHESSComputing/golib/materialscommons"
)

// MCProvider represents Material Commons provider
type MCProvider struct {
	ProjectName string
	Verbose     int
}

// Init function initializes MaterialsCommons publisher
func (m *MCProvider) Init() {
	if m.ProjectName == "" {
		m.ProjectName = srvConfig.Config.MaterialsCommons.ProjectName
	}
	m.Verbose = srvConfig.Config.MaterialsCommons.Verbose
}

// Publish provides publication of dataset with did and description
func (m *MCProvider) Publish(authors []string, did, description string, record map[string]any, publish bool) (string, string, error) {
	doi, doiLink, err := materialscommons.Publish(authors, did, m.ProjectName, description, record, publish, m.Verbose)
	if err != nil {
		return doi, doiLink, fmt.Errorf("[golib.doi.MCProvider.Publish] materialscommons.Publish error: %w", err)
	}
	return doi, doiLink, nil
}

// MakePublic provides publication of draft DOI
func (m *MCProvider) MakePublic(doi string) error {
	return materialscommons.MakePublic(doi, m.Verbose)
}
