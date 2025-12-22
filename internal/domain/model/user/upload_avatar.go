package xuser

import "mime/multipart"

// ===== REQUEST từ HTTP =====

type UploadAvatarLocalRequest struct {
	UserID uint
	File   *multipart.FileHeader `form:"file" validate:"required,image"`
}

type UploadAvatarCloudRequest struct {
	UserID uint
	File   *multipart.FileHeader `form:"file" validate:"required,image"`
}

// ===== PAYLOAD cho QUEUE =====

type UploadAvatarLocalQueuePayload struct {
	UserID uint
	File   *multipart.FileHeader
}

type UploadAvatarCloudQueuePayload struct {
	UserID    uint
	File      *multipart.FileHeader
	OldAvatar string
}

// ===== INPUT từ JOB callback =====

type UploadAvatarLocalInput struct {
	UserID   uint
	Filename string
	Filepath string
}

type UploadAvatarCloudInput struct {
	UserID    uint
	PublicID  string
	SecureURL string
}

// ===== DELETE CLOUD =====

type DeleteCloudAssetPayload struct {
	PublicID string
}
