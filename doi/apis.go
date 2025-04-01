package doi

// Provider represents generic DOI interface
type Provider interface {
	Publish(did, description string, record map[string]any, publish bool) (string, string, error)
	MakePublic(doi string) error
}
