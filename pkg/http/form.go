package xhttp

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
)

var (
	ErrInvalidFormType     = errors.New("request content type is not multipart/form-data")
	ErrRequiredField       = errors.New("required field is missing")
	ErrFileIndexOutOfRange = errors.New("file index out of range")
)

type MultipartForm struct {
	Form     *multipart.Form
	MaxSize  int64
	FormData map[string][]string
	Files    map[string][]*multipart.FileHeader
}

func NewMultipartForm(r *http.Request, maxSize int64) (*MultipartForm, error) {
	// Check content type
	contentType := r.Header.Get("Content-Type")
	if contentType == "" || contentType[:len("multipart/form-data")] != "multipart/form-data" {
		return nil, ErrInvalidFormType
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(maxSize); err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	if r.MultipartForm == nil {
		return nil, fmt.Errorf("multipart form is nil")
	}

	return &MultipartForm{
		Form:     r.MultipartForm,
		MaxSize:  maxSize,
		FormData: r.MultipartForm.Value,
		Files:    r.MultipartForm.File,
	}, nil
}

func (f *MultipartForm) GetString(key string) (string, bool) {
	if values, exists := f.FormData[key]; exists && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func (f *MultipartForm) GetStringOrDefault(key, defaultValue string) string {
	if value, ok := f.GetString(key); ok {
		return value
	}
	return defaultValue
}

// GetStringRequired returns a string value or an error if not found
func (f *MultipartForm) GetStringRequired(key string) (string, error) {
	if val, exists := f.GetString(key); exists {
		return val, nil
	}
	return "", fmt.Errorf("%w: %s", ErrRequiredField, key)
}

// GetBool returns a boolean value from the form
func (f *MultipartForm) GetBool(key string) (bool, bool) {
	if val, exists := f.GetString(key); exists {
		boolVal, err := strconv.ParseBool(val)
		if err == nil {
			return boolVal, true
		}
	}
	return false, false
}

// GetBoolOrDefault returns a boolean value or a default value if not found
func (f *MultipartForm) GetBoolOrDefault(key string, defaultValue bool) bool {
	if val, exists := f.GetBool(key); exists {
		return val
	}
	return defaultValue
}

// GetInt returns an integer value from the form
func (f *MultipartForm) GetInt(key string) (int, bool) {
	if val, exists := f.GetString(key); exists {
		intVal, err := strconv.Atoi(val)
		if err == nil {
			return intVal, true
		}
	}
	return 0, false
}

// GetIntOrDefault returns an integer value or a default value if not found
func (f *MultipartForm) GetIntOrDefault(key string, defaultValue int) int {
	if val, exists := f.GetInt(key); exists {
		return val
	}
	return defaultValue
}

// GetIntRequired returns an integer value or an error if not found
func (f *MultipartForm) GetIntRequired(key string) (int, error) {
	if val, exists := f.GetInt(key); exists {
		return val, nil
	}
	return 0, fmt.Errorf("%w: %s", ErrRequiredField, key)
}

// GetUint returns an unsigned integer value from the form
func (f *MultipartForm) GetUint(key string) (uint, bool) {
	if val, exists := f.GetString(key); exists {
		intVal, err := strconv.ParseUint(val, 10, 32)
		if err == nil {
			return uint(intVal), true
		}
	}
	return 0, false
}

// GetUintOrDefault returns an unsigned integer value or a default value if not found
func (f *MultipartForm) GetUintOrDefault(key string, defaultValue uint) uint {
	if val, exists := f.GetUint(key); exists {
		return val
	}
	return defaultValue
}

// GetUintRequired returns an unsigned integer value or an error if not found
func (f *MultipartForm) GetUintRequired(key string) (uint, error) {
	if val, exists := f.GetUint(key); exists {
		return val, nil
	}
	return 0, fmt.Errorf("%w: %s", ErrRequiredField, key)
}

// GetFiles returns all file headers for a given key
func (f *MultipartForm) GetFiles(key string) ([]*multipart.FileHeader, bool) {
	if files, exists := f.Files[key]; exists && len(files) > 0 {
		return files, true
	}
	return nil, false
}

// GetFile returns the first file header for a given key
func (f *MultipartForm) GetFile(key string) (*multipart.FileHeader, bool) {
	if files, exists := f.GetFiles(key); exists {
		return files[0], true
	}
	return nil, false
}

// ValidateFileIndex checks if a file index is valid
func (f *MultipartForm) ValidateFileIndex(fileField string, fileIndex int) (*multipart.FileHeader, error) {
	files, exists := f.GetFiles(fileField)
	if !exists {
		return nil, fmt.Errorf("no files found with key %s", fileField)
	}

	if fileIndex < 0 || fileIndex >= len(files) {
		return nil, fmt.Errorf("%w: index %d, file count %d", ErrFileIndexOutOfRange, fileIndex, len(files))
	}

	return files[fileIndex], nil
}
