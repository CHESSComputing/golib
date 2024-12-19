package doi

type Provider interface {
	Init()
	Publish(did, description string) (string, string, error)
}
