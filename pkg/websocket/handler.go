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

func (h *Handler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		websocket.Handler(h.Server.Handle).ServeHTTP(
			c.Response(),
			c.Request(),
		)
		return nil
	}
}
