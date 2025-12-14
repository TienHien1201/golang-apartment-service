package xutils

import (
	"encoding/json"
)

// AfterScanArray converts a string (JSON array or comma-separated) to string slice
func AfterScanArray(input string) []string {
	if input == "" {
		return []string{}
	}

	// Try to parse as JSON array first
	var result []string
	if err := json.Unmarshal([]byte(input), &result); err == nil {
		return result
	}

	// If not JSON, use StringToArray for comma-separated handling
	return StringToArray(input)
}
