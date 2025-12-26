package jobs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"thomas.vn/apartment_service/internal/domain/consts"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	user2 "thomas.vn/apartment_service/internal/domain/usecase"
	xfile "thomas.vn/apartment_service/pkg/file"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xqueue "thomas.vn/apartment_service/pkg/queue"
)

type UploadUserAvatarJob struct {
	logger      *xlogger.Logger
	fileService *xfile.HTTPFile
	userUC      user2.UserUsecase
}

func NewUploadUserAvatarJob(
	logger *xlogger.Logger,
	fileService *xfile.HTTPFile,
	userUC user2.UserUsecase,
) *UploadUserAvatarJob {
	return &UploadUserAvatarJob{
		logger:      logger,
		fileService: fileService,
		userUC:      userUC,
	}
}

func (j *UploadUserAvatarJob) Name() string {
	return consts.UploadUserAvatarJobName
}

func (j *UploadUserAvatarJob) Type() xqueue.MessageType {
	return consts.UploadUserAvatarJobType
}

func (j *UploadUserAvatarJob) Handle(ctx context.Context, payload interface{}) error {
	req, ok := payload.(*xuser.UploadAvatarLocalQueuePayload)
	if !ok {
		return fmt.Errorf("invalid payload")
	}

	dir := "attachments/images/avatar"

	fullPath, err := j.fileService.Upload(req.File, dir)
	if err != nil {
		return err
	}

	oldName := filepath.Base(fullPath)
	newName := "local-" + oldName

	oldPath := fullPath
	newPath := filepath.Join(dir, newName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("rename avatar file failed: %w", err)
	}

	return j.userUC.ProcessUploadLocal(ctx, &xuser.UploadAvatarLocalInput{
		UserID:   req.UserID,
		Filename: newName,
		Filepath: newPath,
	})
}
