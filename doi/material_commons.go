package doi

import (
	materialCommons "github.com/CHESSComputing/golib/MaterialCommons"
	srvConfig "github.com/CHESSComputing/golib/config"
)

// MaterialCommons represents Material Commons Publisher
type MaterialCommons struct {
	Name string
}

// Init function initializes MaterialCommons publisher
func (m *MaterialCommons) Init() {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	m.Name = srvConfig.Config.MaterialCommons.ProjectName
}

// Publish provides publication of dataset with did and description
func (m *MaterialCommons) Publish(did, description string) (string, string, error) {
	doi, doiLink, err := materialCommons.Publish(did, description)
	return doi, doiLink, err
}
