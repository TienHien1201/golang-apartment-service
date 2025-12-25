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
	DeletedBy    int        `json:"deleted_by"`
	IsDeleted    int        `json:"is_deleted"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type ListChatMessageRequest struct {
	query.PaginationOptions
	query.DateRangeOptions
	query.SortOptions

	ChatGroupID int `json:"-" query:"-"`
}

type CreateChatMessageRequest struct {
	ChatGroupID  int    `json:"chat_group_id"`
	UserIDSender int    `json:"user_id_sender"`
	MessageText  string `json:"message_text"`
}

type Response struct {
	ID          int       `json:"id"`
	MessageText string    `json:"message_text"`
	CreatedAt   time.Time `json:"created_at"`
	ChatGroupID int       `json:"chat_group_id"`
	Sender      Sender    `json:"sender"`
}

type Sender struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	RoleID   int    `json:"role_id"`
}

type Row struct {
	ID          int
	ChatGroupID int
	MessageText string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	UserID   int
	FullName string
	Avatar   string
	RoleID   int
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
