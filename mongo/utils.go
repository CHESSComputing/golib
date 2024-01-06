package mongo

import (
	"sort"

	"github.com/CHESSComputing/golib/utils"
)

// MapKeys helper function to return keys from a map
func MapKeys(rec map[string]any) []string {
	keys := make([]string, 0, len(rec))
	for k := range rec {
		keys = append(keys, k)
	}
	sort.Sort(utils.StringList(keys))
	return keys
}
