package utils

// utils module
//
// Copyright (c) 2024 - Valentin Kuznetsov <vkuznet AT gmail dot com>
//

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	constraints "golang.org/x/exp/constraints"
)

// FindInList helper function to find item in a list
func FindInList(a string, arr []string) bool {
	for _, e := range arr {
		if e == a {
			return true
		}
	}
	return false
}

// InList helper function to check item in a list
func InList[T ListEntry](a T, list []T) bool {
	check := 0
	for _, b := range list {
		if b == a {
			check++
		}
	}
	if check != 0 {
		return true
	}
	return false
}

// MapKeys returns string keys from a map
func MapKeys(rec map[string]interface{}) []string {
	keys := make([]string, 0, len(rec))
	for k := range rec {
		keys = append(keys, k)
	}
	sort.Sort(StringList(keys))
	return keys
}

// MapIntKeys returns int keys from a map
func MapIntKeys(rec map[int]interface{}) []int {
	keys := make([]int, 0, len(rec))
	for k := range rec {
		keys = append(keys, k)
	}
	return keys
}

// EqualLists helper function to compare list of strings
func EqualLists(list1, list2 []string) bool {
	count := 0
	for _, k := range list1 {
		if InList(k, list2) {
			count++
		} else {
			return false
		}
	}
	if len(list2) == count {
		return true
	}
	return false
}

// CheckEntries helper function to check that entries from list1 are all appear in list2
func CheckEntries(list1, list2 []string) bool {
	var out []string
	for _, k := range list1 {
		if InList(k, list2) {
			//             count += 1
			out = append(out, k)
		}
	}
	if len(out) == len(list1) {
		return true
	}
	return false
}

// List2Set helper function to convert input list into set
func List2Set[T ListEntry](arr []T) []T {
	var out []T
	for _, key := range arr {
		if !InList(key, out) {
			out = append(out, key)
		}
	}
	return out
}

// IsInt helper function to test if given value is integer
func IsInt(val string) bool {
	return PatternInt.MatchString(val)
}

// IsFloat helper function to test if given value is a float
func IsFloat(val string) bool {
	return PatternFloat.MatchString(val)
}

// Sum helper function to perform sum operation over provided array of values
func Sum(data []interface{}) float64 {
	out := 0.0
	for _, val := range data {
		if val != nil {
			//             out += val.(float64)
			switch v := val.(type) {
			case float64:
				out += v
			case json.Number:
				vv, e := v.Float64()
				if e == nil {
					out += vv
				}
			case int64:
				out += float64(v)
			}
		}
	}
	return out
}

// Max helper function to perform Max operation over provided array of values
func Max(data []interface{}) float64 {
	out := 0.0
	for _, val := range data {
		if val != nil {
			switch v := val.(type) {
			case float64:
				if v > out {
					out = v
				}
			case json.Number:
				vv, e := v.Float64()
				if e == nil && vv > out {
					out = vv
				}
			case int64:
				if float64(v) > out {
					out = float64(v)
				}
			}
		}
	}
	return out
}

// Min helper function to perform Min operation over provided array of values
func Min(data []interface{}) float64 {
	out := float64(^uint(0) >> 1) // largest int
	for _, val := range data {
		if val == nil {
			continue
		}
		switch v := val.(type) {
		case float64:
			if v < out {
				out = v
			}
		case json.Number:
			vv, e := v.Float64()
			if e == nil && vv < out {
				out = vv
			}
		case int64:
			if float64(v) < out {
				out = float64(v)
			}
		}
	}
	return out
}

// IntList implement sort for []int type
type IntList []int

// Len provides length of the []int type
func (s IntList) Len() int { return len(s) }

// Swap implements swap function for []int type
func (s IntList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements less function for []int type
func (s IntList) Less(i, j int) bool { return s[i] < s[j] }

// Int64List implement sort for []int type
type Int64List []int64

// Len provides length of the []int64 type
func (s Int64List) Len() int { return len(s) }

// Swap implements swap function for []int64 type
func (s Int64List) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements less function for []int64 type
func (s Int64List) Less(i, j int) bool { return s[i] < s[j] }

// StringList implement sort for []string type
type StringList []string

// Len provides length of the []int type
func (s StringList) Len() int { return len(s) }

// Swap implements swap function for []int type
func (s StringList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements less function for []int type
func (s StringList) Less(i, j int) bool { return s[i] < s[j] }

// ListEntry identifies types used by list's generics function
type ListEntry interface {
	int | int64 | float64 | string
}

// Set converts input list into set
func Set[T ListEntry](arr []T) []T {
	var out []T
	for _, v := range arr {
		if !InList(v, out) {
			out = append(out, v)
		}
	}
	return out
}

// sortSlice helper function on any ordered generic list
// https://gosamples.dev/generics-sort-slice/
func sortSlice[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

// OrderedSet implementa ordered set
func OrderedSet[T ListEntry](list []T) []T {
	out := Set(list)
	sortSlice(out)
	return out
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal[T ListEntry](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// ListFiles lists files in a given directory
func ListFiles(dir string) []string {
	var out []string
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	for _, f := range entries {
		if !f.IsDir() {
			out = append(out, f.Name())
		}
	}
	return out
}

// Insert inserts value into array at zero position
func Insert(arr []interface{}, val interface{}) []interface{} {
	arr = append(arr, val)
	copy(arr[1:], arr[0:])
	arr[0] = val
	return arr
}

// UpdateOrderedDict returns new ordered list from given ordered dicts
func UpdateOrderedDict(omap, nmap map[int][]string) map[int][]string {
	for idx, list := range nmap {
		if entries, ok := omap[idx]; ok {
			entries = append(entries, list...)
			omap[idx] = entries
		} else {
			omap[idx] = list
		}
	}
	return omap
}

// UniqueFormValues returns unique list of values from http.Request.FormValue
func UniqueFormValues(vals []string) []string {
	vals = List2Set(vals)
	// the url forms provide values in []string form
	// loop below remote duplicates and separate possible item values by empty space
	var items []string
	for _, v := range vals {
		sarr := strings.Split(v, " ")
		sarr = List2Set(sarr)
		for _, s := range sarr {
			items = append(items, s)
		}
	}
	return OrderedSet[string](items)
}
