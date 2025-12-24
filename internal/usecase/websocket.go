package usecase

import (
	"context"
	"strconv"

	"thomas.vn/apartment_service/internal/domain/model/chatgroup"
	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type ChatUcase struct {
	logger        *xlogger.Logger
	chatGroupUC   usecase.ChatGroupUsecase
	chatMessageUC usecase.ChatMessageUsecase
}

func NewChatUcase(
	logger *xlogger.Logger,
	chatGroupUC usecase.ChatGroupUsecase,
	chatMessageUC usecase.ChatMessageUsecase,
) *ChatUcase {
	return &ChatUcase{
		logger:        logger,
		chatGroupUC:   chatGroupUC,
		chatMessageUC: chatMessageUC,
	}
}
func (u *ChatUcase) CreateRoom(ctx context.Context, req *chatgroup.CreateChatGroupRequest) (string, error) {

	group, err := u.chatGroupUC.CreateChatGroup(ctx, req)
	if err != nil {
		u.logger.Error("CreateRoom failed", xlogger.Error(err))
		return "", err
	}

	return strconv.FormatInt(group.ID, 10), nil
}

func (u *ChatUcase) SendMessage(ctx context.Context, request *chatmessage.CreateChatMessageRequest) error {

	if _, err := u.chatMessageUC.SendMessage(ctx, request); err != nil {
		u.logger.Error("SendMessage failed", xlogger.Error(err))
		return err
	}

	return nil
}
