package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// writeWait is the maximum time allowed to write a message.
	writeWait = 10 * time.Second
	// pongWait is the maximum time to wait for a pong response.
	pongWait = 60 * time.Second
	// pingPeriod is how often to send pings (must be less than pongWait).
	pingPeriod = (pongWait * 9) / 10 // 54 seconds
	// maxMessageSize is the maximum size of an inbound message (8 KB).
	maxMessageSize = 8 * 1024
)

// Client represents a single authenticated WebSocket connection.
type Client struct {
	conn   *websocket.Conn
	UserID string
	send   chan []byte

	// rooms tracks which chat rooms this client has joined (for cleanup on disconnect).
	rooms map[int]struct{}
	mu    sync.Mutex

	// limiter enforces a per-client message rate limit.
	limiter *rateLimiter
}

func newClient(conn *websocket.Conn, userID string) *Client {
	return &Client{
		conn:    conn,
		UserID:  userID,
		send:    make(chan []byte, 512),
		rooms:   make(map[int]struct{}),
		limiter: newRateLimiter(10, time.Second), // 10 messages/second max
	}
}

func (c *Client) JoinRoom(room int) {
	c.mu.Lock()
	c.rooms[room] = struct{}{}
	c.mu.Unlock()
}

func (c *Client) LeaveRoom(room int) {
	c.mu.Lock()
	delete(c.rooms, room)
	c.mu.Unlock()
}

func (c *Client) InRoom(room int) bool {
	c.mu.Lock()
	_, ok := c.rooms[room]
	c.mu.Unlock()
	return ok
}

// AllRooms returns a snapshot of the rooms this client has joined.
func (c *Client) AllRooms() []int {
	c.mu.Lock()
	defer c.mu.Unlock()
	rooms := make([]int, 0, len(c.rooms))
	for r := range c.rooms {
		rooms = append(rooms, r)
	}
	return rooms
}

// WritePump reads outgoing messages from c.send and forwards them to the
// WebSocket connection. It also sends periodic pings to keep the connection alive.
// Must be run in its own goroutine. Closes the connection when the channel is closed.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel — send a clean close frame.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
