package doi

type Provider interface {
	Init()
	Publish(did, description string, record any, publish bool) (string, string, error)
}
