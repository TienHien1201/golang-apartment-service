package repository

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model/chatgroup"
)

type ChatGroupRepository interface {
	ListChatGroupsWithMembers(ctx context.Context, req *chatgroup.ListChatGroupRequest) ([]*chatgroup.ListResponse, int64, error)
	CreateChatGroup(ctx context.Context, chatGroup *chatgroup.ChatGroup) (*chatgroup.ChatGroup, error)
	AddMembers(ctx context.Context, req *chatgroup.CreateMemberRequest) error
}
