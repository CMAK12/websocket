package websocket

import (
	"encoding/json"

	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
)

type Client struct {
	ID   string
	send chan []byte
	conn *websocket.Conn
	hub  *Hub
}

type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
	To      string `json:"to"`
}

func NewClient(
	ID string,
	conn *websocket.Conn,
	hub *Hub,
) *Client {
	return &Client{
		ID:   ID,
		send: make(chan []byte, 256),
		conn: conn,
		hub:  hub,
	}
}

func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			c.hub.logger.Error("Failed to unmarshal message", zap.Error(err))
			continue
		}

		msg.From = c.ID
		c.hub.Broadcast(msg)
	}
}

func (c *Client) Write() {
	defer c.conn.Close()

	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}
