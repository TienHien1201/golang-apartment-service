package ws

import (
	"context"
	"encoding/json"
	"fmt"
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
		Rooms: map[int]bool{},
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

// ================= ROUTER =================

func (s *Server) dispatch(c *Client, msg Message) {
	fmt.Printf("🟢 Received WS type=%s, payload=%s\n", msg.Type, string(msg.Payload))
	switch msg.Type {
	case "AUTH":
		s.handleAuth(c, msg)
	case "JOIN_ROOM":
		s.handleJoinRoom(c, msg)
	case "CREATE_ROOM":
		s.handleCreateRoom(c, msg)
	case "SEND_MESSAGE":
		s.handleSendMessage(c, msg)
	}
}

// ================= HANDLERS =================

func (s *Server) handleAuth(c *Client, msg Message) {
	var p struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return
	}

	claims, err := s.Token.VerifyAccessToken(p.AccessToken)
	if err != nil {
		return
	}

	c.UserID = strconv.Itoa(int(claims.UserID))
}

func (s *Server) handleJoinRoom(c *Client, msg Message) {
	var p struct {
		ChatGroupID int `json:"chatGroupId"`
	}
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return
	}

	room := p.ChatGroupID
	s.Hub.Join(room, c)
	c.Rooms[room] = true
}

func (s *Server) handleCreateRoom(c *Client, msg Message) {
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

	s.Hub.Join(roomID, c)
	c.Rooms[roomID] = true

	resp, _ := json.Marshal(map[string]any{
		"type": "CREATE_ROOM",
		"data": map[string]any{
			"chatGroupId": roomID,
		},
	})

	c.Send <- resp
}
func (s *Server) handleSendMessage(c *Client, msg Message) {
	fmt.Println("📥 handleSendMessage called with payload:", string(msg.Payload))
	var p struct {
		Message     string `json:"message"`
		ChatGroupID int    `json:"chatGroupId"`
		AccessToken string `json:"accessToken"`
	}
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return
	}

	claims, err := s.Token.VerifyAccessToken(p.AccessToken)
	if err != nil {
		return
	}

	// 👉 ensure client joined room
	room := p.ChatGroupID
	if !c.Rooms[room] {
		s.Hub.Join(room, c)
		c.Rooms[room] = true
	}

	req := &chatmessage.CreateChatMessageRequest{
		ChatGroupID:  p.ChatGroupID,
		UserIDSender: int(claims.UserID),
		MessageText:  p.Message,
	}

	if err := s.ChatUC.SendMessage(context.Background(), req); err != nil {
		return
	}

	event, _ := json.Marshal(map[string]any{
		"type": "SEND_MESSAGE",
		"data": map[string]any{
			"message_text": p.Message,
			"userIdSender": claims.UserID,
			"chatGroupId":  p.ChatGroupID,
		},
	})

	s.Hub.Broadcast(room, event)
}

// ================= WRITE =================

func (s *Server) write(c *Client) {
	for msg := range c.Send {
		if err := websocket.Message.Send(c.Conn, msg); err != nil {
			return
		}
	}
}
