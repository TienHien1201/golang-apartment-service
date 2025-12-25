package usecase

import (
	"context"

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
func (u *ChatUcase) CreateRoom(ctx context.Context, req *chatgroup.CreateChatGroupRequest) (int, error) {

	group, err := u.chatGroupUC.CreateChatGroup(ctx, req)
	if err != nil {
		u.logger.Error("CreateRoom failed", xlogger.Error(err))
		return 0, err
	}

	return group.ID, nil
}

func (u *ChatUcase) SendMessage(ctx context.Context, req *chatmessage.CreateChatMessageRequest) (*chatmessage.Response, error) {

	resp, err := u.chatMessageUC.SendMessage(ctx, req)
	if err != nil {
		u.logger.Error("SendMessage failed", xlogger.Error(err))
		return nil, err
	}

	return resp, nil
}
