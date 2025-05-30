package doi

import (
	srvConfig "github.com/CHESSComputing/golib/config"
	materialscommons "github.com/CHESSComputing/golib/materialscommons"
)

// MCProvider represents Material Commons provider
type MCProvider struct {
	Name    string
	Verbose int
}

// Init function initializes MaterialsCommons publisher
func (m *MCProvider) Init() {
	m.Name = srvConfig.Config.MaterialsCommons.ProjectName
}

// Publish provides publication of dataset with did and description
func (m *MCProvider) Publish(did, description string, record map[string]any, publish bool) (string, string, error) {
	doi, doiLink, err := materialscommons.Publish(did, description, record, publish)
	return doi, doiLink, err
}

// MakePublic provides publication of draft DOI
func (m *MCProvider) MakePublic(doi string) error {
	return materialscommons.MakePublic(doi)
}
