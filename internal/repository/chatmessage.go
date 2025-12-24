package repository

import (
	"context"

	"gorm.io/gorm"
	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
	"thomas.vn/apartment_service/internal/domain/repository"
	xlogger "thomas.vn/apartment_service/pkg/logger"
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

func (r *chatMessageRepository) ListChatMessages(ctx context.Context, req *chatmessage.ListChatMessageRequest) ([]*chatmessage.ChatMessage, int64, error) {
	var chatMessages []*chatmessage.ChatMessage
	var total int64

	query := r.chatMessageTable.
		WithContext(ctx).
		Model(&chatmessage.ChatMessage{})

	// Apply filters
	if req.IsDeleted != 0 {
		query = query.Where("is_deleted = ?", req.IsDeleted)
	}
	if req.FromDate != "" {
		query = query.Where(req.RangeBy+" >= ?", req.FromDate+" 00:00:00")
	}
	if req.ToDate != "" {
		query = query.Where(req.RangeBy+" <= ?", req.ToDate+" 23:59:59")
	}

	// Get total count if not exclude
	if !req.ExcludeTotal {
		if err := query.Count(&total).Error; err != nil {
			r.logger.Error("Count chatmessage failed", xlogger.Error(err))
			return nil, 0, err
		}
	}

	// Apply pagination
	if req.Page > 0 && req.Limit > 0 {
		query = query.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)
	}

	// Apply sorting
	if req.SortBy != "" && req.OrderBy != "" {
		query = query.Order(req.SortBy + " " + req.OrderBy)
	}

	// Execute query
	if err := query.Find(&chatMessages).Error; err != nil {
		r.logger.Error("List chatmessage failed", xlogger.Error(err))
		return nil, 0, err
	}

	return chatMessages, total, nil
}
