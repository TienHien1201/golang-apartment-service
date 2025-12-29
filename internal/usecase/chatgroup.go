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

func (u *ChatGroupUsecase) CreateChatGroup(ctx context.Context, req *chatgroup.CreateChatGroupRequest) (*chatgroup.ChatGroup, error) {

	if len(req.TargetUserIDs) == 0 {
		return nil, fmt.Errorf("target users is required")
	}

	userIDs := append(req.TargetUserIDs, req.OwnerID)
	userIDs = uniqueInt64(userIDs)

	entity := &chatgroup.ChatGroup{
		Name:      req.Name,
		OwnerID:   req.OwnerID,
		IsDeleted: consts.NotDeleted,
	}

	// ================= CHAT 1–1 =================
	if len(userIDs) == 2 {
		keyObj, err := chatgroup.BuildChatOneKey(userIDs)
		if err != nil {
			return nil, err
		}

		entity.KeyForChatOne = keyObj
	}

	createdGroup, err := u.chatGroupRepositoty.CreateChatGroup(ctx, entity)
	if err != nil {
		u.logger.Error("CreateChatGroup failed", xlogger.Error(err))
		return nil, err
	}

	memberReq := &chatgroup.CreateMemberRequest{
		ChatGroupID: int64(createdGroup.ID),
		UserIDs:     userIDs,
	}

	if err := u.chatGroupRepositoty.AddMembers(ctx, memberReq); err != nil {
		u.logger.Error("AddMembers failed", xlogger.Error(err))
		return nil, err
	}

	return createdGroup, nil
}

func uniqueInt64(input []int64) []int64 {
	m := make(map[int64]struct{}, len(input))
	out := make([]int64, 0, len(input))
	for _, v := range input {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
