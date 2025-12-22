package jobs

import (
	"context"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"thomas.vn/apartment_service/internal/domain/consts"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xqueue "thomas.vn/apartment_service/pkg/queue"
)

type DeleteCloudinaryAssetJob struct {
	logger *xlogger.Logger
	cld    *cloudinary.Cloudinary
}

func NewDeleteCloudinaryAssetJob(
	logger *xlogger.Logger,
	cld *cloudinary.Cloudinary,
) *DeleteCloudinaryAssetJob {
	return &DeleteCloudinaryAssetJob{
		logger: logger,
		cld:    cld,
	}
}

func (j *DeleteCloudinaryAssetJob) Name() string {
	return "delete_cloudinary_asset_job"
}

func (j *DeleteCloudinaryAssetJob) Type() xqueue.MessageType {
	return consts.DeleteCloudinaryAssetJobType
}
func (j *DeleteCloudinaryAssetJob) Handle(
	ctx context.Context,
	payload interface{},
) error {

	req, ok := payload.(*xuser.DeleteCloudAssetPayload)
	if !ok {
		return xhttp.NewAppError(
			"ERR_INVALID_PAYLOAD",
			"cloudinary",
			"invalid payload",
			http.StatusBadRequest,
		)
	}

	_, err := j.cld.Upload.Destroy(
		ctx,
		uploader.DestroyParams{
			PublicID:     req.PublicID,
			ResourceType: "image",
		},
	)

	if err != nil {
		j.logger.Error(
			"delete cloudinary asset failed",
			xlogger.Error(err),
			xlogger.String("public_id", req.PublicID),
		)
		return err
	}

	j.logger.Info(
		"cloudinary asset deleted",
		xlogger.String("public_id", req.PublicID),
	)

	return nil
}
