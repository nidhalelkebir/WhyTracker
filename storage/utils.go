package storage

import (
	"encoding/json"
)

// stringSliceToJSON converts a string slice to JSON string
func stringSliceToJSON(slice []string) string {
	if len(slice) == 0 {
		return "[]"
	}
	data, err := json.Marshal(slice)
	if err != nil {
		return "[]"
	}
	return string(data)
}

// jsonToStringSlice converts JSON string to string slice
func jsonToStringSlice(jsonStr string) []string {
	if jsonStr == "" || jsonStr == "[]" {
		return []string{}
	}
	var slice []string
	if err := json.Unmarshal([]byte(jsonStr), &slice); err != nil {
		return []string{}
	}
	return slice
}
