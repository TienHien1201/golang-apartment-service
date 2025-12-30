package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type chatMessageUsecase struct {
	logger                *xlogger.Logger
	chatMessageRepository repository.ChatMessageRepository
}

func NewChatMessageUsecase(logger *xlogger.Logger, chatMessageRepository repository.ChatMessageRepository) usecase.ChatMessageUsecase {
	return &chatMessageUsecase{
		logger:                logger,
		chatMessageRepository: chatMessageRepository,
	}
}
func (u *chatMessageUsecase) ListChatMessages(ctx context.Context, req *chatmessage.ListChatMessageRequest) ([]*chatmessage.Response, int64, error) {

	return u.chatMessageRepository.ListChatMessages(ctx, req)
}

func (u *chatMessageUsecase) SendMessage(ctx context.Context, req *chatmessage.CreateChatMessageRequest) (*chatmessage.Response, error) {
	entity := &chatmessage.ChatMessage{
		ChatGroupID:  req.ChatGroupID,
		UserIDSender: req.UserIDSender,
		MessageText:  req.MessageText,
	}

	row, err := u.chatMessageRepository.CreateChatMessage(ctx, entity)
	if err != nil {
		u.logger.Error("SendMessage failed", xlogger.Error(err))
		return nil, err
	}

	resp := &chatmessage.Response{
		ID:          row.ID,
		MessageText: row.MessageText,
		CreatedAt:   row.CreatedAt,
		ChatGroupID: row.ChatGroupID,
	}

	resp.Sender.ID = row.UserID
	resp.Sender.FullName = row.FullName
	resp.Sender.Avatar = row.Avatar
	resp.Sender.RoleID = row.RoleID

	return resp, nil
}
