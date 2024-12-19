package doi

type Provider interface {
	Init()
	Publish(did, description string, record any) (string, string, error)
}
