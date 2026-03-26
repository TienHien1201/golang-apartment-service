package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// AllowOrigins controls CORS for WebSocket upgrades.
	// In production, replace with a strict allowlist.
	CheckOrigin: func(_ *http.Request) bool { return true },
}

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
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Router /ws [get]
func (h *Handler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		h.Server.Handle(conn)
		return nil
	}
}
