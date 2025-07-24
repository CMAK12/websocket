package websocket

import (
	"sync"

	"go.uber.org/zap"
)

type Hub struct {
	clients    map[*Client]bool
	mux        sync.RWMutex
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	logger     *zap.Logger
}

func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)

		case message := <-h.broadcast:
			h.routeMessage(message)
		}
	}
}

func (h *Hub) handleRegister(client *Client) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.clients[client] = true

	h.logger.Info("Client registered")
}

func (h *Hub) handleUnregister(client *Client) {
	h.mux.Lock()
	defer h.mux.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		h.logger.Info("Client unregistered")
	}
}

func (h *Hub) routeMessage(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
			h.logger.Warn("Client send channel closed")
		}
	}
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}
