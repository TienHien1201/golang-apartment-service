package model

import "time"

type ChatGroupMembers struct {
	ID          int        `json:"id"`
	UserID      int64      `json:"user_id"`
	ChatGroupID int64      `json:"chat_group_id"`
	DeletedBy   int        `json:"deleted_by"`
	IsDeleted   int        `json:"is_deleted"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (ChatGroupMembers) TableName() string {
	return "chat_group_members"
}
