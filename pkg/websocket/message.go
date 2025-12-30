package ws

import "encoding/json"

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
type Auth struct {
	AccessToken string `json:"accessToken" example:"string"`
}
type JoinGroupPayload struct {
	ChatGroupID int `json:"chatGroupId" example:"123"`
}
type SendMessagePayload struct {
	ChatGroupID int    `json:"chatGroupId" example:"123"`
	Message     string `json:"message" example:"Hello"`
	AccessToken string `json:"accessToken" example:"string"`
}

type CreateRoomPayload struct {
	Name          string  `json:"name" example:"Room name"`
	TargetUserIDs []int64 `json:"targetUserIDs" example:"1,2,3"`
	AccessToken   string  `json:"accessToken" example:"string"`
}
