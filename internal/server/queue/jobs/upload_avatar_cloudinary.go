package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"thomas.vn/apartment_service/internal/domain/consts"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	user2 "thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xqueue "thomas.vn/apartment_service/pkg/queue"
)

type UploadAvatarCloudJob struct {
	logger *xlogger.Logger
	cld    *cloudinary.Cloudinary
	userUC user2.UserUsecase
}

func NewUploadAvatarCloudJob(
	logger *xlogger.Logger,
	cld *cloudinary.Cloudinary,
	userUC user2.UserUsecase,
) *UploadAvatarCloudJob {
	return &UploadAvatarCloudJob{
		logger: logger,
		cld:    cld,
		userUC: userUC,
	}
}

func (j *UploadAvatarCloudJob) Name() string {
	return consts.UploadAvatarCloudJobName
}

func (j *UploadAvatarCloudJob) Type() xqueue.MessageType {
	return consts.UploadAvatarCloudJobType
}

func (j *UploadAvatarCloudJob) Handle(ctx context.Context, payload interface{}) error {
	req, ok := payload.(*xuser.UploadAvatarCloudQueuePayload)
	if !ok {
		return xhttp.NewAppError(
			"ERR_INVALID_PAYLOAD",
			"avatar",
			"invalid payload",
			http.StatusBadRequest,
		)
	}

	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	uploadResult, err := j.cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			Folder:       "images/avatar",
			ResourceType: "image",
			PublicID:     fmt.Sprintf("user_%d", req.UserID),
			Overwrite:    api.Bool(true),
		},
	)
	if err != nil {
		j.logger.Error(
			"Cloudinary upload avatar failed",
			xlogger.Error(err),
			xlogger.Uint("user_id", req.UserID),
		)
		return err
	}

	return j.userUC.ProcessUploadCloud(ctx, &xuser.UploadAvatarCloudInput{
		UserID:    req.UserID,
		PublicID:  uploadResult.PublicID,
		SecureURL: uploadResult.SecureURL,
	})
}
