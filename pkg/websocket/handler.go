package ws

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type Handler struct {
	Server *Server
}

func NewHandler(server *Server) *Handler {
	return &Handler{Server: server}
}

// @Summary WebSocket server
// @Description Connect to WebSocket server to handle chat events. All payloads are JSON.
// @Tags websocket
// @Accept json
// @Produce json
//
//	@Param AUTH body ws.Auth true "Authenticate user" example({
//	  "type": "AUTH",
//	  "payload": {
//	    "accessToken": "string"
//	  }
//	})
//
//	@Param JOIN_ROOM body ws.JoinGroupPayload true "Join a chat room" example({
//	  "type": "JOIN_ROOM",
//	  "payload": {
//	    "chatGroupId": 123
//	  }
//	})
//
//	@Param CREATE_ROOM body ws.CreateRoomPayload true "Create a chat room" example({
//	  "type": "CREATE_ROOM",
//	  "payload": {
//	    "name": "Room name",
//	    "targetUserIDs": [1,2,3],
//	    "accessToken": "string"
//	  }
//	})
//
//	@Param SEND_MESSAGE body ws.SendMessagePayload true "Send message to a chat room" example({
//	  "type": "SEND_MESSAGE",
//	  "payload": {
//	    "chatGroupId": 123,
//	    "message": "Hello",
//	    "accessToken": "string"
//	  }
//	})
//
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Router /ws [get]
func (h *Handler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		websocket.Handler(h.Server.Handle).ServeHTTP(
			c.Response(),
			c.Request(),
		)
		return nil
	}
}
