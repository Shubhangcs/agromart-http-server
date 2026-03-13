package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket user.
type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub keeps track of all active WebSocket connections, keyed by userID.
// Multiple tabs / devices for the same user are supported via a slice.
type Hub struct {
	mu      sync.RWMutex
	clients map[string][]*Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string][]*Client),
	}
}

// Register adds a client connection for a user.
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.UserID] = append(h.clients[c.UserID], c)
}

// Unregister removes a specific client connection and closes its Send channel.
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns := h.clients[c.UserID]
	remaining := conns[:0]
	for _, existing := range conns {
		if existing != c {
			remaining = append(remaining, existing)
		}
	}
	if len(remaining) == 0 {
		delete(h.clients, c.UserID)
	} else {
		h.clients[c.UserID] = remaining
	}
	close(c.Send)
}

// Send delivers a raw JSON message to all connections belonging to userID.
// Returns true if the user had at least one active connection.
func (h *Hub) Deliver(userID string, msg []byte) bool {
	h.mu.RLock()
	conns := h.clients[userID]
	h.mu.RUnlock()

	if len(conns) == 0 {
		return false
	}
	for _, c := range conns {
		select {
		case c.Send <- msg:
		default:
			// Slow consumer — drop the message rather than blocking.
		}
	}
	return true
}

// IsOnline reports whether a user has at least one open connection.
func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID]) > 0
}
