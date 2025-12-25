package repository

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
)

type ChatMessageRepository interface {
	ListChatMessages(ctx context.Context, req *chatmessage.ListChatMessageRequest) ([]*chatmessage.Response, int64, error)
	CreateChatMessage(ctx context.Context, chatMessage *chatmessage.ChatMessage) (*chatmessage.Row, error)
}
