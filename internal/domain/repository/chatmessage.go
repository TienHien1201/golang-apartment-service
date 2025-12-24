package repository

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
)

type ChatMessageRepository interface {
	ListChatMessages(
		ctx context.Context,
		req *chatmessage.ListChatMessageRequest,
	) ([]*chatmessage.ChatMessage, int64, error)
}
