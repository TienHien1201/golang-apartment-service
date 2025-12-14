package xfile

import "errors"

const (
	PhotoMaxWidth  = 1400
	PhotoMaxHeight = 1400
	DefaultMaxSize = 2 << 20 // 2MB
	DefaultMinSize = 1 << 10 // 1KB
)

var (
	ErrInvalidFileExtension = errors.New("invalid file extension")
	ErrFileTooLarge         = errors.New("file size too large")
	ErrFileTooSmall         = errors.New("file size too small")
	ErrImageDimension       = errors.New("invalid image dimensions")
	ErrFileExists           = errors.New("file already exists")
	ErrInvalidMimeType      = errors.New("invalid mime type")
)
