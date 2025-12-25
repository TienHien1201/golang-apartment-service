package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
)

type ChatMessageUsecase interface {
	ListChatMessages(ctx context.Context, req *chatmessage.ListChatMessageRequest) ([]*chatmessage.Response, int64, error)
	SendMessage(ctx context.Context, req *chatmessage.CreateChatMessageRequest) (*chatmessage.Response, error)
}
