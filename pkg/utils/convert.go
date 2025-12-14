package xutils

import (
	"fmt"
	"strconv"
	"strings"
)

// IntToInt8 converts *int to *int8 safely
func IntToInt8(val *int) *int8 {
	if val == nil {
		return nil
	}
	v := int8(*val)
	return &v
}

// IntToInt16 converts *int to *int16 safely
func IntToInt16(val *int) *int16 {
	if val == nil {
		return nil
	}
	v := int16(*val)
	return &v
}

// IntToString converts *int to *string safely
func IntToString(val *int) *string {
	if val == nil {
		return nil
	}
	v := strconv.Itoa(*val)
	return &v
}

// ConvertFilePathToJSONArray converts a comma-separated string of file paths to a JSON array
func ConvertFilePathToJSONArray(filePath string) string {
	// Handle empty input
	if filePath == "" {
		return "[]"
	}

	// Remove any existing quotes and brackets
	cleanPath := strings.Trim(filePath, "[]\"")
	// Split by comma if multiple paths
	paths := strings.Split(cleanPath, ",")

	// Convert each path to JSON string format
	jsonPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		// Trim spaces
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		// Replace backslashes with forward slashes
		path = strings.ReplaceAll(path, "\\", "/")
		// Remove leading/trailing slashes for consistency
		path = strings.Trim(path, "/")
		// Add single leading forward slash
		path = "/" + path
		// Escape forward slashes for JSON
		path = strings.ReplaceAll(path, "/", "\\/")
		jsonPaths = append(jsonPaths, fmt.Sprintf("\"%s\"", path))
	}

	// Join all paths with commas and wrap in square brackets
	return fmt.Sprintf("[%s]", strings.Join(jsonPaths, ","))
}

// ConvertFilePathsToJSONArray converts a slice of file paths to a JSON array
func ConvertFilePathsToJSONArray(paths []string) string {
	if len(paths) == 0 {
		return "[]"
	}

	// Convert each path to JSON string format
	jsonPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		// Trim spaces
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		// Replace backslashes with forward slashes
		path = strings.ReplaceAll(path, "\\", "/")
		// Remove leading/trailing slashes for consistency
		path = strings.Trim(path, "/")
		// Add single leading forward slash
		path = "/" + path
		// Escape forward slashes for JSON
		path = strings.ReplaceAll(path, "/", "\\/")
		jsonPaths = append(jsonPaths, fmt.Sprintf("\"%s\"", path))
	}

	// Join all paths with commas and wrap in square brackets
	return fmt.Sprintf("[%s]", strings.Join(jsonPaths, ","))
}
