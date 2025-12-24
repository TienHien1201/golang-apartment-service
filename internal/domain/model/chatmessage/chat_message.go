package chatmessage

import (
	"time"

	"thomas.vn/apartment_service/pkg/query"
)

type ChatMessage struct {
	ID           int        `json:"id"`
	ChatGroupID  int        `json:"chat_group_id"`
	UserIDSender int        `json:"user_id_sender"`
	MessageText  string     `json:"message_text"`
	DeleteBy     string     `json:"delete_by"`
	IsDeleted    int        `json:"is_deleted"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type ListChatMessageRequest struct {
	query.PaginationOptions
	query.DateRangeOptions
	query.SortOptions
	IsDeleted int `query:"is_deleted" validate:"omitempty,oneof=0 1"`
}

type ListChatMessageResponse struct {
	Rows  []*ChatMessage `json:"rows"`
	Total int64          `json:"total,omitempty"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
