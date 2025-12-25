package repository

import (
	"context"

	"gorm.io/gorm"
	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
	"thomas.vn/apartment_service/internal/domain/repository"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xutils "thomas.vn/apartment_service/pkg/utils"
)

type chatMessageRepository struct {
	logger           *xlogger.Logger
	chatMessageTable *gorm.DB
}

func NewChatMessageRepository(logger *xlogger.Logger, db *gorm.DB) repository.ChatMessageRepository {
	return &chatMessageRepository{
		logger:           logger,
		chatMessageTable: db.Table("chat_messages"),
	}

}
func (r *chatMessageRepository) ListChatMessages(
	ctx context.Context,
	req *chatmessage.ListChatMessageRequest,
) ([]*chatmessage.Response, int64, error) {

	var rows []*chatmessage.Row
	var total int64

	db := r.chatMessageTable.WithContext(ctx).
		Table("chat_messages cm").
		Joins("JOIN users u ON u.id = cm.user_id_sender").
		Select(`
			cm.id,
			cm.chat_group_id,
			cm.message_text,
			cm.created_at,
			u.id AS user_id,
			u.full_name,
			u.avatar,
			u.role_id
		`)

	if req.ChatGroupID != 0 {
		db = db.Where("cm.chat_group_id = ?", req.ChatGroupID)
	}

	if !req.ExcludeTotal {
		db.Count(&total)
	}

	if req.Page > 0 && req.Limit > 0 {
		db = db.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)
	}

	if req.SortBy != "" && req.OrderBy != "" {
		db = db.Order("cm." + req.SortBy + " " + req.OrderBy)
	}

	if err := db.Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	res := make([]*chatmessage.Response, 0, len(rows))
	for _, r := range rows {
		res = append(res, &chatmessage.Response{
			ID:          r.ID,
			ChatGroupID: r.ChatGroupID,
			MessageText: r.MessageText,
			CreatedAt:   r.CreatedAt,
			Sender: chatmessage.Sender{
				ID:       r.UserID,
				FullName: r.FullName,
				Avatar:   r.Avatar,
				RoleID:   r.RoleID,
			},
		})
	}

	return res, total, nil
}

func (r *chatMessageRepository) CreateChatMessage(ctx context.Context, msg *chatmessage.ChatMessage) (*chatmessage.Row, error) {

	now := xutils.GetTimeNow()
	msg.CreatedAt = now
	msg.UpdatedAt = now

	if err := r.chatMessageTable.WithContext(ctx).Create(msg).Error; err != nil {
		r.logger.Error("CreateChatMessage failed", xlogger.Error(err))
		return nil, err
	}

	var row chatmessage.Row
	err := r.chatMessageTable.WithContext(ctx).
		Table("chat_messages cm").
		Joins("JOIN users u ON u.id = cm.user_id_sender").
		Select(`
			cm.id,
			cm.chat_group_id,
			cm.message_text,
			cm.created_at,
			u.id AS user_id,
			u.full_name,
			u.avatar,
			u.role_id
		`).
		Where("cm.id = ?", msg.ID).
		Scan(&row).Error

	if err != nil {
		return nil, err
	}

	return &row, nil
}
