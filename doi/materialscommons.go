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
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	m.Name = srvConfig.Config.MaterialsCommons.ProjectName
}

// Publish provides publication of dataset with did and description
func (m *MCProvider) Publish(did, description string, record any, publish bool) (string, string, error) {
	doi, doiLink, err := materialscommons.Publish(did, description, record, publish)
	return doi, doiLink, err
}
