package service

import (
	"sync"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
)

type Publisher interface {
	Publish(userID int, notify dto.Notification)
}

type Client struct {
	Ch chan dto.Notification
}

type Hub struct {
	mu      sync.RWMutex
	clients map[int]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[int]map[*Client]struct{})}
}

func (h *Hub) Register(userID int) *Client {
	c := &Client{Ch: make(chan dto.Notification, 16)}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*Client]struct{})
	}
	h.clients[userID][c] = struct{}{}
	return c
}

func (h *Hub) Unregister(userID int, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	set, ok := h.clients[userID]
	if !ok {
		return
	}
	if _, exists := set[c]; !exists {
		return
	}
	delete(set, c)
	close(c.Ch)
	if len(set) == 0 {
		delete(h.clients, userID)
	}
}

func (h *Hub) Publish(userID int, n dto.Notification) {
	var slow []*Client

	h.mu.RLock()
	for c := range h.clients[userID] {
		select {
		case c.Ch <- n:
		default:
			slow = append(slow, c)
		}
	}
	h.mu.RUnlock()

	for _, c := range slow {
		h.Unregister(userID, c)
	}
}
