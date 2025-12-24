package usecase

import (
	"context"
	"fmt"

	"thomas.vn/apartment_service/internal/domain/consts"
	"thomas.vn/apartment_service/internal/domain/model/chatgroup"
	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type ChatGroupUsecase struct {
	logger              *xlogger.Logger
	chatGroupRepositoty repository.ChatGroupRepository
}

func NewChatGroupUsecase(logger *xlogger.Logger, chatGroupRepositoty repository.ChatGroupRepository) usecase.ChatGroupUsecase {
	return &ChatGroupUsecase{
		logger:              logger,
		chatGroupRepositoty: chatGroupRepositoty,
	}
}

func (u *ChatGroupUsecase) ListChatGroups(
	ctx context.Context,
	req *chatgroup.ListChatGroupRequest,
) ([]*chatgroup.ListResponse, int64, error) {

	return u.chatGroupRepositoty.ListChatGroupsWithMembers(ctx, req)
}

func (u *ChatGroupUsecase) CreateChatGroup(
	ctx context.Context,
	req *chatgroup.CreateChatGroupRequest,
) (*chatgroup.ChatGroup, error) {

	if len(req.TargetUserIDs) == 0 {
		return nil, fmt.Errorf("target users is required")
	}

	entity := &chatgroup.ChatGroup{Name: req.Name, OwnerID: req.OwnerID, IsDeleted: consts.NotDeleted}

	createdGroup, err := u.chatGroupRepositoty.CreateChatGroup(ctx, entity)
	if err != nil {
		u.logger.Error("CreateChatGroup failed", xlogger.Error(err))
		return nil, err
	}

	memberReq := &chatgroup.CreateMemberRequest{
		ChatGroupID: createdGroup.ID,
		UserIDs:     append(req.TargetUserIDs, req.OwnerID),
	}

	if err := u.chatGroupRepositoty.AddMembers(ctx, memberReq); err != nil {
		u.logger.Error("AddMembers failed", xlogger.Error(err))
		return nil, err
	}

	return createdGroup, nil
}
