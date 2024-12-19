package doi

import "github.com/CHESSComputing/golib/services"

var _httpReadRequest, _httpWriteRequest *services.HttpRequest

type Provider interface {
	Init()
	Publish(did, description string) (string, string, error)
}
