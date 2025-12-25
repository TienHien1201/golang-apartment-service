package ws

import (
	"sync"

	"golang.org/x/net/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID string
	Send   chan []byte
	Rooms  map[int]bool
	Mu     sync.Mutex
}
