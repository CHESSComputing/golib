package doi

type Provider interface {
	Publish(did, description string) (string, string, error)
}
