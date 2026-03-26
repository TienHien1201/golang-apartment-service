package ws

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gorilla/websocket"

	"thomas.vn/apartment_service/internal/domain/model/chatgroup"
	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
	"thomas.vn/apartment_service/internal/domain/usecase"
)

type Server struct {
	Hub    *Hub
	ChatUC usecase.ChatUsecase
	Token  usecase.TokenUsecase
}

// Handle upgrades the HTTP connection, creates a Client, starts its pumps,
// and blocks until the connection is closed.
func (s *Server) Handle(conn *websocket.Conn) {
	client := newClient(conn, "")

	go client.WritePump()
	s.readPump(client)
}

// readPump reads messages from the WebSocket connection in a loop.
// It cleans up the client from all rooms on exit.
func (s *Server) readPump(c *Client) {
	defer func() {
		c.mu.Lock()
		for room := range c.rooms {
			s.Hub.Leave(room, c)
		}
		c.mu.Unlock()
		close(c.send)
	}()

	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		s.dispatch(c, msg)
	}
}

// ================= ROUTER =================

func (s *Server) dispatch(c *Client, msg Message) {
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
	var p Auth
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
	var p JoinGroupPayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return
	}

	room := p.ChatGroupID
	s.Hub.Join(room, c)
	c.JoinRoom(room)
}

func (s *Server) handleCreateRoom(c *Client, msg Message) {
	var p CreateRoomPayload
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
	c.JoinRoom(roomID)

	resp, _ := json.Marshal(map[string]any{
		"type": "CREATE_ROOM",
		"data": map[string]any{
			"chatGroupId": roomID,
		},
	})

	select {
	case c.send <- resp:
	default:
	}
}

func (s *Server) handleSendMessage(c *Client, msg Message) {
	var p SendMessagePayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return
	}

	claims, err := s.Token.VerifyAccessToken(p.AccessToken)
	if err != nil {
		return
	}

	room := p.ChatGroupID
	if !c.InRoom(room) {
		s.Hub.Join(room, c)
		c.JoinRoom(room)
	}

	req := &chatmessage.CreateChatMessageRequest{
		ChatGroupID:  p.ChatGroupID,
		UserIDSender: int(claims.UserID),
		MessageText:  p.Message,
	}

	resp, err := s.ChatUC.SendMessage(context.Background(), req)
	if err != nil {
		return
	}

	event, _ := json.Marshal(map[string]any{
		"type": "SEND_MESSAGE",
		"data": resp,
	})

	s.Hub.Broadcast(room, event)
}
