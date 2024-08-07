package ql

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"testing"
)

// TestServiceMap
func TestServiceMap(t *testing.T) {
	srv1 := "service1"
	srvKeys1 := []string{"foo", "bla"}
	srv2 := "service2"
	srvKeys2 := []string{"abc", "xyz"}
	smap := make(map[string][]string)
	smap[srv1] = srvKeys1
	smap[srv2] = srvKeys2
	services := []string{srv1, srv2}

	file, err := os.CreateTemp("", "srvmap-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	var qlRecords []QLRecord
	qlRecords = append(qlRecords, QLRecord{Service: "service1", Key: "foo", DataType: "string"})
	qlRecords = append(qlRecords, QLRecord{Service: "service2", Key: "abc", DataType: "string"})
	data, err := json.Marshal(qlRecords)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write(data); err != nil {
		t.Fatal(err)
	}
	file.Close()
	var qlMgr QLManager
	err = qlMgr.Init(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(srvKeys1)
	sort.Strings(srvKeys2)
	sort.Strings(services)
	if !reflect.DeepEqual(qlMgr.Keys(srv1), []string{"foo"}) {
		//     if !reflect.DeepEqual(qlMgr.Keys(srv1), srvKeys1) {
		t.Errorf("service %s, wrong keys %v != %v", srv1, qlMgr.Keys(srv1), srvKeys1)
	}
	if !reflect.DeepEqual(qlMgr.Keys(srv2), []string{"abc"}) {
		//     if !reflect.DeepEqual(qlMgr.Keys(srv2), srvKeys2) {
		t.Errorf("service %s, wrong keys %v != %v", srv2, qlMgr.Keys(srv2), srvKeys2)
	}
	if !reflect.DeepEqual(qlMgr.Services(), services) {
		t.Errorf("wrong services %v != %v", qlMgr.Services(), services)
	}

	// test services queries
	query := "foo:1 abc:[1,2,3]"
	sqMap, err := qlMgr.ServiceQueries(query)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("services queries: %+v", sqMap)
}
