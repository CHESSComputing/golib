package utils

import (
	"fmt"
	"sort"
	"strings"
)

// DIDKeys provide sorted, lower-case list of did keys from comma separated list of attributes
func DIDKeys(attrs string) []string {
	if attrs == "" {
		attrs = "btr,beamline,cycle,sample"
	}
	attrs = strings.Replace(attrs, " ", "", -1)
	var keys []string
	for _, k := range strings.Split(attrs, ",") {
		if k != "" {
			keys = append(keys, strings.ToLower(k))
		}
	}
	sort.Strings(keys)
	return keys
}

func CreateDID(rec map[string]any, attrs, sep, div string) string {
	didKeys := DIDKeys(attrs)
	var did string
	for _, k := range didKeys {
		v, ok := rec[k]
		if !ok {
			continue
		}
		key := strings.ToLower(k)
		var val string
		switch vvv := v.(type) {
		case []string:
			val = strings.ToLower(fmt.Sprintf("%v", strings.Join(vvv, ",")))
		case []int:
			var arr []string
			for _, i := range vvv {
				arr = append(arr, fmt.Sprintf("%d", i))
			}
			val = strings.ToLower(fmt.Sprintf("%s", strings.Join(arr, ",")))
		case []int64:
			var arr []string
			for _, i := range vvv {
				arr = append(arr, fmt.Sprintf("%d", i))
			}
			val = strings.ToLower(fmt.Sprintf("%s", strings.Join(arr, ",")))
		case []any:
			var arr []string
			for _, i := range vvv {
				arr = append(arr, fmt.Sprintf("%v", i))
			}
			val = strings.ToLower(fmt.Sprintf("%v", strings.Join(arr, ",")))
		default:
			val = strings.ToLower(fmt.Sprintf("%v", vvv))
		}
		if InList(key, didKeys) {
			did = fmt.Sprintf("%s%s%s%s%v", did, sep, key, div, val)
		}
	}
	return did
}
