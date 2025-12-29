package chatgroup

import (
	"time"

	"thomas.vn/apartment_service/pkg/query"
)

type ChatGroup struct {
	ID            int         `json:"id"`
	KeyForChatOne *ChatOneKey `json:"key_for_chat_one"`
	Name          string      `json:"name"`
	OwnerID       int64       `json:"owner_id"`
	DeletedBy     int         `json:"deleted_by"`
	IsDeleted     int         `json:"is_deleted"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     *time.Time  `json:"deleted_at"`
}

type ListChatGroupRequest struct {
	query.PaginationOptions
	query.DateRangeOptions
	query.SortOptions
	IsDeleted int  `query:"is_deleted" validate:"omitempty,oneof=0 1"`
	IsOne     bool `query:"isOne"`
}

type CreateChatGroupRequest struct {
	Name          string  `json:"name"`
	OwnerID       int64   `json:"-"`
	TargetUserIDs []int64 `json:"target_user_ids"`
}
type CreateMemberRequest struct {
	ChatGroupID int64
	UserIDs     []int64
}
type ListResponse struct {
	ID               int64            `json:"id"`
	Name             string           `json:"name"`
	ChatGroupMembers []MemberResponse `json:"ChatGroupMembers"`
}

type MemberResponse struct {
	UserID int64        `json:"user_id"`
	Users  UserResponse `json:"Users"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
}
type Row struct {
	GroupID   int64
	GroupName string
	UserID    int64
	FullName  string
	Avatar    string `gorm:"column:avatar"`
}

func (ChatGroup) TableName() string {
	return "chat_groups"
}
