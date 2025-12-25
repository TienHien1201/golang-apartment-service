package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model/chatgroup"
	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
)

type ChatUsecase interface {
	CreateRoom(ctx context.Context, request *chatgroup.CreateChatGroupRequest) (int, error)
	SendMessage(ctx context.Context, request *chatmessage.CreateChatMessageRequest) error
}
