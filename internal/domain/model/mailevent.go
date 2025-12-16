package model

import xqueue "thomas.vn/apartment_service/pkg/queue"

type MailPayload struct {
	Type     xqueue.MessageType `json:"type"`
	Email    string             `json:"email"`
	FullName string             `json:"full_name,omitempty"`
}
