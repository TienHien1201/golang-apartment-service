// Package fileadapter provides an infrastructure adapter that bridges
// pkg/file.HTTPFile (the concrete file implementation) with the domain's
// service.FileService interface.
//
// Clean Architecture rule: the domain layer (service.FileService) may NOT
// import infrastructure packages. This adapter sits in the infrastructure
// layer and performs the type conversion, keeping the domain pure.
package fileadapter

import (
	"mime/multipart"

	"thomas.vn/apartment_service/internal/domain/service"
	xfile "thomas.vn/apartment_service/pkg/file"
)

type adapter struct {
	impl *xfile.HTTPFile
}

// New wraps an *xfile.HTTPFile and returns a service.FileService that
// the domain and usecase layers can depend on without knowing about pkg/file.
func New(impl *xfile.HTTPFile) service.FileService {
	return &adapter{impl: impl}
}

func (a *adapter) Download(fileURL string) (service.File, error) {
	f, err := a.impl.Download(fileURL)
	if err != nil {
		return service.File{}, err
	}
	return service.File{
		Content:     f.Content,
		FileName:    f.FileName,
		ContentType: f.ContentType,
		Size:        f.Size,
	}, nil
}

func (a *adapter) GetFileType(fileName string) string {
	return a.impl.GetFileType(fileName)
}

func (a *adapter) Upload(fileHeader *multipart.FileHeader, dstPath string) (string, error) {
	return a.impl.Upload(fileHeader, dstPath)
}

func (a *adapter) Delete(filePath string) error {
	return a.impl.Delete(filePath)
}

func (a *adapter) CopyFile(srcPath string, dstPath string) (string, error) {
	return a.impl.CopyFile(srcPath, dstPath)
}
