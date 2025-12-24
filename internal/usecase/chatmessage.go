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

func (u *chatMessageUsecase) ListChatMessages(ctx context.Context, req *chatmessage.ListChatMessageRequest) ([]*chatmessage.ChatMessage, int64, error) {
	chatMessages, total, err := u.chatMessageRepository.ListChatMessages(ctx, req)
	if err != nil {
		u.logger.Error("Failed to list Chat Messages", xlogger.Error(err))
		return nil, 0, err
	}

	return chatMessages, total, nil
}
