package ql

import (
	"testing"
)

// TestParseQuery
func TestParseQuery(t *testing.T) {
	Verbose = 1
	query := "bla:1 foo:2"
	spec, err := ParseQuery(query)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("input query %s spec=%+v\n", query, spec)
	for _, k := range []string{"bla", "foo"} {
		if _, ok := spec[k]; !ok {
			t.Errorf("unexpected key %s found\n", k)
		}
	}

	// test 2: use MongoDB QL query
	querySpec := `{"bla":1, "foo": 2}`
	spec, err = ParseQuery(querySpec)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("input query %s spec=%+v\n", querySpec, spec)
}
