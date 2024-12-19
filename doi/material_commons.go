package doi

import (
	materialCommons "github.com/CHESSComputing/golib/MaterialCommons"
)

type MaterialCommons struct {
}

func (m *MaterialCommons) Publish(did, description string) (string, string, error) {
	doi, doiLink, err := materialCommons.Publish(did, description)
	return doi, doiLink, err
}
