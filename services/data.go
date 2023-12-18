package server

import "github.com/CHESSComputing/golib/mongo"

// MetaRecord represents meta-data record used for injection
type MetaRecord struct {
	Schema string
	Record mongo.Record
}
