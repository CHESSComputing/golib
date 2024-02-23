package utils

import (
	"testing"
)

// TestDIDKeys
func TestDIDKeys(t *testing.T) {
	keys := DIDKeys(", a, b, c, ")
	if keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
		t.Error("Fail TestDIDKeys")
	}
}

// TestCreateDID
func TestCreateDID(t *testing.T) {
	rec := make(map[string]any)
	rec["foo"] = 1
	rec["bla"] = "value"
	rec["arr"] = []int{1, 2, 3}
	attrs := "bla,foo,arr"
	sep := "/"
	div := ":"
	did := CreateDID(rec, attrs, sep, div)
	expect := "/arr:1,2,3/bla:value/foo:1"
	if did != expect {
		t.Errorf("Fail TestCreateDID did=%s, expect=%s\n", did, expect)
	}
}
