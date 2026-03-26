package service

import "mime/multipart"

// File represents a downloaded or processed file's content and metadata.
// This is the domain's own type — infrastructure adapters convert between
// this type and any framework-specific file types (e.g. pkg/file.File).
type File struct {
	Content     []byte
	FileName    string
	ContentType string
	Size        int64
}

// FileService defines file I/O operations needed by the domain layer.
// The concrete implementation (pkg/file.HTTPFile) is bridged via
// internal/infrastructure/fileadapter to keep the domain import-free.
type FileService interface {
	Download(fileURL string) (File, error)
	GetFileType(fileName string) string
	Upload(fileHeader *multipart.FileHeader, dstPath string) (string, error)
	Delete(filePath string) error
	CopyFile(srcPath string, dstPath string) (string, error)
}
