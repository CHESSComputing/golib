package doi

import (
	materialCommons "github.com/CHESSComputing/golib/MaterialCommons"
	srvConfig "github.com/CHESSComputing/golib/config"
)

// MCProvider represents Material Commons provider
type MCProvider struct {
	Name string
}

// Init function initializes MaterialCommons publisher
func (m *MCProvider) Init() {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	m.Name = srvConfig.Config.MaterialCommons.ProjectName
}

// Publish provides publication of dataset with did and description
func (m *MCProvider) Publish(did, description string, record any) (string, string, error) {
	doi, doiLink, err := materialCommons.Publish(did, description, record)
	return doi, doiLink, err
}
