package usecase

import (
	"context"
	"fmt"
	"strings"

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

func (u *chatMessageUsecase) SendMessage(ctx context.Context, req *chatmessage.CreateChatMessageRequest) (*chatmessage.ChatMessage, error) {

	if req.ChatGroupID == 0 {
		return nil, fmt.Errorf("chat_group_id is required")
	}
	if req.UserIDSender == 0 {
		return nil, fmt.Errorf("user_id_sender is required")
	}
	if strings.TrimSpace(req.MessageText) == "" {
		return nil, fmt.Errorf("message_text is required")
	}

	entity := &chatmessage.ChatMessage{
		ChatGroupID:  req.ChatGroupID,
		UserIDSender: req.UserIDSender,
		MessageText:  req.MessageText,
	}

	createdMsg, err := u.chatMessageRepository.CreateChatMessage(ctx, entity)
	if err != nil {
		u.logger.Error("SendMessage failed", xlogger.Error(err))
		return nil, err
	}

	return createdMsg, nil
}
