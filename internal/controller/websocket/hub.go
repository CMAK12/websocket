package websocket

import (
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

type Hub struct {
	clients    map[string]*Client
	mux        sync.RWMutex
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	logger     *zap.Logger
}

func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan Message, 256),
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

	h.clients[client.ID] = client

	h.logger.Info("Client registered")
}

func (h *Hub) handleUnregister(client *Client) {
	h.mux.Lock()
	defer h.mux.Unlock()

	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients, client.ID)
		close(client.send)
		h.logger.Info("Client unregistered")
	}
}

func (h *Hub) routeMessage(message Message) {
	h.mux.RLock()
	defer h.mux.RUnlock()

	receiver := h.clients[message.To]

	if receiver != nil {
		data, err := json.Marshal(message)
		if err != nil {
			h.logger.Error("Failed to marshal message", zap.Error(err))
			return
		}

		receiver.send <- data
	}
}

func (h *Hub) Broadcast(message Message) {
	h.broadcast <- message
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}
