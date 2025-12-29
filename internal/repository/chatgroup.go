package repository

import (
	"context"
	"errors"
	"sort"

	"gorm.io/gorm"
	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/model/chatgroup"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xutils "thomas.vn/apartment_service/pkg/utils"
)

type ChatGroupRepository struct {
	logger               *xlogger.Logger
	chatGroupTable       *gorm.DB
	chatGroupMemberTable *gorm.DB
}

func NewChatGroupRepository(logger *xlogger.Logger, db *gorm.DB) *ChatGroupRepository {
	return &ChatGroupRepository{
		logger:               logger,
		chatGroupTable:       db.Table("chat_groups"),
		chatGroupMemberTable: db.Table("chat_group_members"),
	}
}

func (r *ChatGroupRepository) ListChatGroupsWithMembers(ctx context.Context, req *chatgroup.ListChatGroupRequest) ([]*chatgroup.ListResponse, int64, error) {

	var rows []*chatgroup.Row
	var total int64

	query := r.chatGroupTable.WithContext(ctx).
		Table("chat_groups cg").
		Select(`
			cg.id   as group_id,
			cg.name as group_name,
			u.id    as user_id,
			u.full_name,
			u.avatar
		`).
		Joins("JOIN chat_group_members cgm ON cgm.chat_group_id = cg.id").
		Joins("JOIN users u ON u.id = cgm.user_id").
		Where("cg.is_deleted = 0")

	if req.IsOne {
		query = query.Where("cg.key_for_chat_one != ''")
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	groupMap := map[int64]*chatgroup.ListResponse{}

	for _, r := range rows {
		if _, ok := groupMap[r.GroupID]; !ok {
			groupMap[r.GroupID] = &chatgroup.ListResponse{
				ID:   r.GroupID,
				Name: r.GroupName,
			}
		}

		groupMap[r.GroupID].ChatGroupMembers = append(
			groupMap[r.GroupID].ChatGroupMembers,
			chatgroup.MemberResponse{
				UserID: r.UserID,
				Users: chatgroup.UserResponse{
					ID:       r.UserID,
					FullName: r.FullName,
					Avatar:   r.Avatar,
				},
			},
		)
	}

	result := make([]*chatgroup.ListResponse, 0, len(groupMap))
	for _, v := range groupMap {
		result = append(result, v)
	}

	total = int64(len(result))
	return result, total, nil
}
func (r *ChatGroupRepository) CreateChatGroup(ctx context.Context, chatGroup *chatgroup.ChatGroup) (*chatgroup.ChatGroup, error) {

	chatGroup.CreatedAt = xutils.GetTimeNow()
	chatGroup.UpdatedAt = xutils.GetTimeNow()

	if err := r.chatGroupTable.WithContext(ctx).Create(chatGroup).Error; err != nil {
		return nil, err
	}

	return chatGroup, nil
}

func (r *ChatGroupRepository) AddMembers(ctx context.Context, req *chatgroup.CreateMemberRequest) error {

	members := make([]*model.ChatGroupMembers, 0, len(req.UserIDs))
	for _, uid := range req.UserIDs {
		members = append(members, &model.ChatGroupMembers{
			ChatGroupID: req.ChatGroupID,
			UserID:      uid,
		})
	}

	if err := r.chatGroupMemberTable.WithContext(ctx).Create(&members).Error; err != nil {
		r.logger.Error("AddMembers failed", xlogger.Error(err))
		return err
	}
	return nil
}

func (r *ChatGroupRepository) FindChatOneByUserIDs(ctx context.Context, userIDs []int64) (*chatgroup.ChatGroup, error) {

	if len(userIDs) != 2 {
		return nil, nil
	}

	sort.Slice(userIDs, func(i, j int) bool {
		return userIDs[i] < userIDs[j]
	})

	var group chatgroup.ChatGroup

	err := r.chatGroupTable.
		WithContext(ctx).
		Where(`
			JSON_CONTAINS(key_for_chat_one, JSON_ARRAY(?), '$.user_ids')
			AND JSON_CONTAINS(key_for_chat_one, JSON_ARRAY(?), '$.user_ids')
		`, userIDs[0], userIDs[1]).
		First(&group).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &group, nil
}
