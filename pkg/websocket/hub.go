package ws

import "sync"

type Hub struct {
	Rooms map[string]map[*Client]bool
	Mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Join(room string, c *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	if h.Rooms[room] == nil {
		h.Rooms[room] = make(map[*Client]bool)
	}
	h.Rooms[room][c] = true
}

func (h *Hub) Leave(room string, c *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	delete(h.Rooms[room], c)
}

func (h *Hub) Broadcast(room string, msg []byte) {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

	for c := range h.Rooms[room] {
		c.Send <- msg
	}
}
