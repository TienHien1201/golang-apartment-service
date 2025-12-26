package xfile

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// UploadOptions defines configuration for file uploads
type UploadOptions struct {
	MaxSize          int64
	MinSize          int64
	AllowedExts      []string
	AllowedMimeTypes []string
	Overwrite        bool
	MinWidth         int
	MinHeight        int
	MaxWidth         int
	MaxHeight        int
	ImageQuality     int
	ScanForViruses   bool
}

// UploadHandler handles file uploads with validation and processing
type UploadHandler struct {
	options  UploadOptions
	fullPath string
}

// NewUploadHandler creates a new upload handler with default options
func NewUploadHandler(options UploadOptions) *UploadHandler {
	if options.MaxSize == 0 {
		options.MaxSize = DefaultMaxSize
	}
	if options.MinSize == 0 {
		options.MinSize = DefaultMinSize
	}
	if options.ImageQuality == 0 {
		options.ImageQuality = 90
	}
	return &UploadHandler{
		options: options,
	}
}

// Upload handles the file upload process
func (u *UploadHandler) Upload(fileHeader *multipart.FileHeader, dstPath string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := u.validate(fileHeader); err != nil {
		return "", err
	}

	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	fileName := generateFileName(fileHeader.Filename, u.options.Overwrite)
	u.fullPath = filepath.Join(dstPath, fileName)

	if err := u.saveFile(file); err != nil {
		return "", err
	}

	if u.isImageFile(fileHeader.Filename) {
		if err := u.processImage(); err != nil {
			os.Remove(u.fullPath) // Cleanup if image processing fails
			return "", err
		}
	}

	if u.options.ScanForViruses {
		if err := u.scanForViruses(); err != nil {
			os.Remove(u.fullPath) // Cleanup if virus scan fails
			return "", err
		}
	}

	return u.fullPath, nil
}

// validate performs all file validations
func (u *UploadHandler) validate(fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !u.isValidExtension(ext) {
		return ErrInvalidFileExtension
	}

	if !u.isValidMimeType(fileHeader.Header.Get("Content-Type")) {
		return ErrInvalidMimeType
	}

	if fileHeader.Size > u.options.MaxSize {
		return ErrFileTooLarge
	}

	if fileHeader.Size < u.options.MinSize {
		return ErrFileTooSmall
	}

	return nil
}

// saveFile saves the uploaded file to disk
func (u *UploadHandler) saveFile(file multipart.File) error {
	out, err := os.Create(u.fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

// processImage handles image-specific processing
func (u *UploadHandler) processImage() error {
	img, err := imaging.Open(u.fullPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if u.options.MinWidth > 0 && width < u.options.MinWidth {
		return fmt.Errorf("%w: min width %dpx required", ErrImageDimension, u.options.MinWidth)
	}

	if u.options.MinHeight > 0 && height < u.options.MinHeight {
		return fmt.Errorf("%w: min height %dpx required", ErrImageDimension, u.options.MinHeight)
	}

	maxWidth := u.options.MaxWidth
	if maxWidth == 0 {
		maxWidth = PhotoMaxWidth
	}

	maxHeight := u.options.MaxHeight
	if maxHeight == 0 {
		maxHeight = PhotoMaxHeight
	}

	if width > maxWidth || height > maxHeight {
		img = imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
	}

	options := imaging.JPEGQuality(u.options.ImageQuality)
	return imaging.Save(img, u.fullPath, options)
}

// scanForViruses performs virus scanning on the uploaded file
func (u *UploadHandler) scanForViruses() error {
	// TODO: Implement virus scanning using a library like clamav
	// For now, this is a placeholder
	return nil
}

// isValidExtension checks if the file extension is allowed
func (u *UploadHandler) isValidExtension(ext string) bool {
	for _, allowed := range u.options.AllowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}

// isValidMimeType checks if the file's MIME type is allowed
func (u *UploadHandler) isValidMimeType(mimeType string) bool {
	for _, allowed := range u.options.AllowedMimeTypes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}

// isImageFile checks if the file is an image based on extension
func (u *UploadHandler) isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case "jpg", "jpeg", "png", "gif":
		return true
	}
	return false
}

// generateFileName creates a unique filename for the uploaded file
func generateFileName(originalName string, overwrite bool) string {
	ext := strings.ToLower(filepath.Ext(originalName)) // .png
	base := strings.TrimSuffix(originalName, ext)

	safeBase := sanitizeFileName(base)

	if overwrite {
		return safeBase + ext
	}

	timestamp := time.Now().Format("20060102150405")
	random := uuid.New().String()
	return fmt.Sprintf("%s_%s_%s%s", safeBase, timestamp, random, ext)
}

// sanitizeFileName cleans the filename to be safe for filesystem use
func sanitizeFileName(name string) string {
	name = strings.TrimSpace(name)

	reg := regexp.MustCompile(`[^\w\s-]`)
	name = reg.ReplaceAllString(name, "")

	regSpace := regexp.MustCompile(`[\s-]+`)
	name = regSpace.ReplaceAllString(name, "-")

	return strings.ToLower(name)
}
