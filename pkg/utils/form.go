package xutils

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	xhttp "thomas.vn/apartment_service/pkg/http"
)

// ParseCompressedData extracts data from form field, decompressing if needed
func ParseCompressedData(form *xhttp.MultipartForm, fieldName string, isCompressed bool) (string, error) {
	// Get data from form
	data, exists := form.GetString(fieldName)
	if !exists || data == "" {
		return "", fmt.Errorf("%s data is required", fieldName)
	}

	// Process based on compression setting
	if isCompressed {
		decompressedBytes, err := Base64DecodeAndDecompress(data)
		if err != nil {
			return "", fmt.Errorf("failed to decompress data: %w", err)
		}
		return string(decompressedBytes), nil
	}

	return data, nil
}

// ParseJSONArray parses a JSON string into a slice of objects
func ParseJSONArray[T any](jsonStr string) ([]T, error) {
	var items []T
	if err := json.Unmarshal([]byte(jsonStr), &items); err != nil {
		return nil, fmt.Errorf("invalid JSON data: %w", err)
	}
	return items, nil
}

// ValidateFileIndexes ensures all file indexes in a slice are valid
func ValidateFileIndexes[T any](items []T, files []*multipart.FileHeader, indexGetter func(T) int) error {
	for i, item := range items {
		fileIndex := indexGetter(item)
		if fileIndex >= 0 && fileIndex >= len(files) {
			return fmt.Errorf("invalid file index in item %d: index %d, available files %d",
				i, fileIndex, len(files))
		}
	}
	return nil
}
