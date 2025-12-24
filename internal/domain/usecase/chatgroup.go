package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model/chatgroup"
)

type ChatGroupUsecase interface {
	ListChatGroups(ctx context.Context, req *chatgroup.ListChatGroupRequest) ([]*chatgroup.ListResponse, int64, error)
	CreateChatGroup(ctx context.Context, req *chatgroup.CreateChatGroupRequest) (*chatgroup.ChatGroup, error)
}
