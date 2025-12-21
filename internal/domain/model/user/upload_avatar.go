package xuser

import "mime/multipart"

type UploadAvatarLocalRequest struct {
	UserID uint
	File   *multipart.FileHeader `form:"file" validate:"required,image"`
}

type UploadAvatarQueuePayload struct {
	UserID uint
	File   *multipart.FileHeader `form:"file" validate:"required,image"`
}

type UploadAvatarLocalInput struct {
	UserID   uint
	Filename string
	Filepath string
}
