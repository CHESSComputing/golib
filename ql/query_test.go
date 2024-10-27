package ql

import (
	"fmt"
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

	// test 3: use regex
	query = `{"did":" /beamline*"}`
	spec, err = ParseQuery(query)
	if err != nil {
		t.Errorf(err.Error())
	}
	if val, ok := spec["did"]; ok {
		vvv := fmt.Sprintf("%v", val)
		if vvv != "map[$regex: /beamline.*]" {
			msg := fmt.Sprintf("parsed query %s does not fit regexp", vvv)
			t.Errorf(msg)
		}
	}
	// test 4: use complex regex query
	query = `{"$or":[{"beamline":".*val.*"},{"btr":".*val.*"}]}`
	spec, err = ParseQuery(query)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("query", query)
	fmt.Println("spec", spec)
	if len(spec) == 0 {
		t.Errorf("empty spec")
	}
}
