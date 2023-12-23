package utils

import (
	"fmt"
	"testing"
)

// TestUtilsInList
func TestUtilsInList(t *testing.T) {
	vals := []string{"1", "2", "3"}
	res := InList("1", vals)
	if res == false {
		t.Error("Fail TestInList")
	}
	res = InList("5", vals)
	if res == true {
		t.Error("Fail TestInList")
	}
}

// TestUtilsSet
func TestUtilsSet(t *testing.T) {
	vals := []string{"a", "b", "c", "a"}
	res := List2Set(vals)
	if len(res) != 3 {
		t.Error("Fail TestUtilsSet")
	}
}

// TestUniqueFormValues
func TestUniqueFormValues(t *testing.T) {
	vals := []string{"a", "b", "a", "b", "a b", "b a"}
	res := UniqueFormValues(vals)
	if len(res) != 2 {
		fmt.Printf("input %+v results %+v", vals, res)
		t.Error("Fail TestUniqueFormValues")
	}
}
