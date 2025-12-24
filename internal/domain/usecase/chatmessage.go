package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
)

type ChatMessageUsecase interface {
	ListChatMessages(ctx context.Context, req *chatmessage.ListChatMessageRequest) ([]*chatmessage.ChatMessage, int64, error)
}
