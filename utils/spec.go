package utils

import (
	"encoding/json"
	"strings"
)

// NormalizeSpec function lower keys in spec
func NormalizeSpec(spec string) string {
	spec = strings.TrimSpace(spec)

	// If it's a JSON object, try to unmarshal and process
	if strings.HasPrefix(spec, "{") && strings.HasSuffix(spec, "}") {
		var raw map[string]any
		if err := json.Unmarshal([]byte(spec), &raw); err != nil {
			return spec // return original if JSON is malformed
		}

		// Create a new map with lowercase keys
		newMap := make(map[string]any)
		for k, v := range raw {
			newMap[strings.ToLower(k)] = v
		}

		// Marshal back to JSON
		newJSON, err := json.Marshal(newMap)
		if err != nil {
			return spec
		}
		return string(newJSON)
	}

	// Handle key:value format
	if idx := strings.Index(spec, ":"); idx != -1 {
		key := strings.ToLower(strings.TrimSpace(spec[:idx]))
		val := strings.TrimSpace(spec[idx+1:])
		return key + ":" + val
	}

	// Otherwise, return as is
	return spec
}
