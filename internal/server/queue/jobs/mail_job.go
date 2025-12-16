package jobs

import (
	"context"
	"fmt"

	"thomas.vn/apartment_service/internal/domain/consts"
	"thomas.vn/apartment_service/internal/domain/model"
	domainUC "thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xqueue "thomas.vn/apartment_service/pkg/queue"
)

type MailJob struct {
	logger *xlogger.Logger
	mailUC domainUC.MailUsecase
}

func NewMailJob(
	logger *xlogger.Logger,
	mailUC domainUC.MailUsecase,
) *MailJob {
	return &MailJob{
		logger: logger,
		mailUC: mailUC,
	}
}

func (j *MailJob) Name() string {
	return consts.MailJobName
}

func (j *MailJob) Type() xqueue.MessageType {
	return consts.MailJobType
}

func (j *MailJob) Handle(ctx context.Context, payload interface{}) error {
	j.logger.Info("Starting MailJob.Handle")
	j.logger.Info(
		"Received payload type",
		xlogger.String("type", fmt.Sprintf("%T", payload)),
	)
	j.logger.Info(
		"Received payload value",
		xlogger.String("value", fmt.Sprintf("%+v", payload)),
	)

	req, ok := payload.(*model.MailPayload)
	if !ok {
		j.logger.Error(
			"Payload type mismatch",
			xlogger.String("expected", "*model.MailPayload"),
			xlogger.String("got", fmt.Sprintf("%T", payload)),
		)
		return fmt.Errorf("invalid mail payload")
	}

	switch req.Type {
	case consts.QueueMailLogin:
		return j.mailUC.SendLoginMail(ctx, req.Email, req.FullName)

	case consts.QueueMailRegister:
		return j.mailUC.SendRegisterMail(ctx, req.Email, req.FullName)

	default:
		j.logger.Error(
			"Unsupported mail type",
			xlogger.String("type", string(req.Type)),
		)
		return fmt.Errorf("unsupported mail type: %s", req.Type)
	}
}
