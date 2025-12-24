package ws

import (
	"context"
	"encoding/json"
	"strconv"

	"golang.org/x/net/websocket"

	"thomas.vn/apartment_service/internal/domain/model/chatgroup"
	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
	"thomas.vn/apartment_service/internal/domain/usecase"
)

type Server struct {
	Hub    *Hub
	ChatUC usecase.ChatUsecase
	Token  usecase.TokenUsecase
}

func (s *Server) Handle(conn *websocket.Conn) {
	defer conn.Close()

	client := &Client{
		Conn:  conn,
		Send:  make(chan []byte, 256),
		Rooms: map[string]bool{},
	}

	go s.write(client)

	for {
		var msg Message
		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			return
		}
		s.dispatch(client, msg)
	}
}

func (s *Server) dispatch(c *Client, msg Message) {
	switch msg.Type {

	// ================= CREATE ROOM =================
	case "CREATE_ROOM":
		var p struct {
			TargetUserIDs []int64 `json:"targetUserIDs"`
			AccessToken   string  `json:"accessToken"`
			Name          string  `json:"name"`
		}
		if err := json.Unmarshal(msg.Payload, &p); err != nil {
			return
		}

		claims, err := s.Token.VerifyAccessToken(p.AccessToken)
		if err != nil {
			return
		}

		req := &chatgroup.CreateChatGroupRequest{
			Name:          p.Name,
			OwnerID:       int64(claims.UserID),
			TargetUserIDs: p.TargetUserIDs,
		}

		roomID, err := s.ChatUC.CreateRoom(context.Background(), req)
		if err != nil {
			return
		}

		room := "chat:" + roomID
		s.Hub.Join(room, c)

		resp, _ := json.Marshal(map[string]any{
			"type": "CREATE_ROOM_SUCCESS",
			"data": map[string]any{
				"chatGroupId": roomID,
			},
		})
		c.Send <- resp

	// ================= SEND MESSAGE =================
	case "SEND_MESSAGE":
		var p struct {
			Message     string `json:"message"`
			ChatGroupID string `json:"chatGroupId"`
			AccessToken string `json:"accessToken"`
		}
		if err := json.Unmarshal(msg.Payload, &p); err != nil {
			return
		}

		claims, err := s.Token.VerifyAccessToken(p.AccessToken)
		if err != nil {
			return
		}

		groupID, err := strconv.ParseInt(p.ChatGroupID, 10, 64)
		if err != nil {
			return
		}

		req := &chatmessage.CreateChatMessageRequest{
			ChatGroupID:  int(groupID),
			UserIDSender: int(claims.UserID),
			MessageText:  p.Message,
		}

		if err := s.ChatUC.SendMessage(context.Background(), req); err != nil {
			return
		}

		event, _ := json.Marshal(map[string]any{
			"type": "SEND_MESSAGE",
			"data": map[string]any{
				"messageText":  p.Message,
				"userIdSender": claims.UserID,
				"chatGroupId":  p.ChatGroupID,
			},
		})

		s.Hub.Broadcast("chat:"+p.ChatGroupID, event)
	}
}

func (s *Server) write(c *Client) {
	for msg := range c.Send {
		if err := websocket.Message.Send(c.Conn, msg); err != nil {
			return
		}
	}
}
