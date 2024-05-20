package ql

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	srvConfig "github.com/CHESSComputing/golib/config"
	utils "github.com/CHESSComputing/golib/utils"
)

// QLKey defines structure of QL key
type QLKey struct {
	Key         string `json:"key"`
	Description string `json:"description,omitempty"`
	Service     string `json:"service"`
	Units       string `json:"units,omitempty"`
	Schema      string `json:"schema,omitempty"`
	DataType    string `json:"type"`
}

// String returns string representation of ql key
func (qlKey *QLKey) String() string {
	// each qmap here is QLKey structure
	desc := qlKey.Description
	if desc == "" {
		desc = "description not available"
	}
	srv := fmt.Sprintf("%s:%s", qlKey.Service, qlKey.Schema)
	if qlKey.Schema == "" {
		srv = qlKey.Service
	}
	key := fmt.Sprintf("%s: (%s) %s", qlKey.Key, srv, desc)
	if qlKey.Units != "" {
		key += fmt.Sprintf(", units:%s", qlKey.Units)
	}
	if qlKey.DataType != "" {
		key += fmt.Sprintf(", data-type:%s", qlKey.DataType)
	}
	return key
}

// QLKeys return list of ql keys
func QLKeys(qlKey string) ([]string, error) {
	var keys []string
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	fname := srvConfig.Config.QL.ServiceMapFile
	file, err := os.Open(fname)
	if err != nil {
		return keys, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return keys, err
	}
	var arr []QLKey
	err = json.Unmarshal(data, &arr)
	if err != nil {
		return keys, err
	}
	var allKeys []string
	for _, elem := range arr {
		if qlKey != "" && elem.Key == qlKey {
			allKeys = append(allKeys, elem.String())
		} else {
			allKeys = append(allKeys, elem.String())
		}
	}
	keys = utils.List2Set[string](allKeys)
	return keys, nil
}
