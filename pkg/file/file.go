package xfile

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"thomas.vn/apartment_service/pkg/retry"
)

type File struct {
	Content     []byte
	FileName    string
	ContentType string
	Size        int64
}

type FileService interface {
	Download(fileURL string) (File, error)
	GetFileType(fileName string) string
	Upload(fileHeader *multipart.FileHeader, dstPath string) (string, error)
	Delete(filePath string) error
	CopyFile(srcPath string, dstPath string) (string, error)
}

type HTTPFile struct {
	httpClient  *http.Client
	retryConfig *retry.Config
}

func NewHTTPFile(client *http.Client) *HTTPFile {
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &HTTPFile{
		httpClient:  client,
		retryConfig: retry.DefaultConfig(),
	}
}

// Download downloads a file from URL and returns it as a File object
func (h *HTTPFile) Download(fileURL string) (File, error) {
	ctx := context.Background()

	downloadFn := func(ctx context.Context) (File, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
		if err != nil {
			// log fileUrl and err
			return File{}, fmt.Errorf("failed to create request for URL %s: %w", fileURL, err)
		}

		resp, err := h.httpClient.Do(req)
		if err != nil {
			return File{}, fmt.Errorf("failed to download file: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return File{}, fmt.Errorf("download failed with status: %d", resp.StatusCode)
		}

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return File{}, fmt.Errorf("failed to read response body: %w", err)
		}

		fileName := h.extractFileName(resp, fileURL)

		return File{
			Content:     content,
			FileName:    fileName,
			ContentType: resp.Header.Get("Content-Type"),
			Size:        int64(len(content)),
		}, nil
	}

	return retry.WithRetry(ctx, h.retryConfig, downloadFn)
}

// Upload handles file upload with validation and processing
// file is pdf doc docx
func (h *HTTPFile) Upload(fileHeader *multipart.FileHeader, dstPath string) (string, error) {
	if fileHeader == nil {
		return "", fmt.Errorf("file header is nil")
	}

	if fileHeader.Size == 0 {
		return "", fmt.Errorf("file is empty")
	}

	// get file type
	fileType := h.GetFileType(fileHeader.Filename)
	ext := filepath.Ext(fileHeader.Filename)

	handler := NewUploadHandler(UploadOptions{
		AllowedMimeTypes: []string{fileType},
		AllowedExts:      []string{ext},
		Overwrite:        false,
		MaxSize:          10 * (1 << 20), // 10MB
		MinSize:          1 << 10,        // 1KB
		ImageQuality:     90,
		ScanForViruses:   false,
		MinWidth:         0,
		MinHeight:        0,
		MaxWidth:         0,
		MaxHeight:        0,
	})

	if !contains(handler.options.AllowedExts, ext) {
		return "", fmt.Errorf("file extension %s is not allowed", ext)
	}

	if !contains(handler.options.AllowedMimeTypes, fileType) {
		return "", fmt.Errorf("file type %s is not allowed", fileType)
	}

	return handler.Upload(fileHeader, dstPath)
}

func (h *HTTPFile) UploadMultiple(fileHeaders []*multipart.FileHeader, dstPath string) ([]string, error) {
	if len(fileHeaders) == 0 {
		return nil, fmt.Errorf("no files to upload")
	}

	// create a slice to store the path of the uploaded files
	filePaths := make([]string, 0, len(fileHeaders))

	// upload each file
	for _, fileHeader := range fileHeaders {
		filePath, err := h.Upload(fileHeader, dstPath)
		if err != nil {
			// if any file upload fails, delete the previously uploaded files
			for _, path := range filePaths {
				_ = h.Delete(path)
			}
			return nil, fmt.Errorf("failed to upload file %s: %w", fileHeader.Filename, err)
		}
		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil
}

// GetFileType returns the MIME type based on file extension
func (h *HTTPFile) GetFileType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".txt":
		return "text/plain"
	case ".csv":
		return "text/csv"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}

// Delete deletes a file from the server
func (h *HTTPFile) Delete(filePath string) error {
	// delete file from server
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// extractFileName tries to get filename from Content-Disposition header
// or falls back to URL path if header is not available
func (h *HTTPFile) extractFileName(resp *http.Response, fileURL string) string {
	// Try to get filename from Content-Disposition header
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		if fileName := extractFileNameFromHeader(contentDisposition); fileName != "" {
			return fileName
		}
	}

	// Fallback to URL path
	return path.Base(fileURL)
}

// extractFileNameFromHeader extracts filename from Content-Disposition header
func extractFileNameFromHeader(header string) string {
	if header == "" {
		return ""
	}

	// Handle both formats:
	// attachment; filename="file.pdf"
	// attachment; filename=file.pdf
	if strings.Contains(header, "filename=") {
		parts := strings.Split(header, "filename=")
		if len(parts) < 2 {
			return ""
		}

		fileName := parts[1]

		// Remove quotes if present
		fileName = strings.Trim(fileName, `"'`)

		// Remove additional parameters if present
		if idx := strings.Index(fileName, ";"); idx > 0 {
			fileName = fileName[:idx]
		}

		return fileName
	}

	return ""
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// CopyFile copies a file from srcPath to dstPath
func (h *HTTPFile) CopyFile(srcPath string, dstPath string) (string, error) {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Download the file from source path
	file, err := h.Download(srcPath)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	// Create a FileHeader from the downloaded file
	header := make(map[string][]string)
	header["Content-Type"] = []string{file.ContentType}

	fileHeader := &multipart.FileHeader{
		Filename: file.FileName,
		Header:   header,
		Size:     file.Size,
	}

	// Upload the file to the destination path
	dstFilePath, err := h.Upload(fileHeader, dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return dstFilePath, nil
}
