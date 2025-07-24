package websocket

import (
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	send chan []byte
	conn *websocket.Conn
	hub  *Hub
}

func NewClient(
	conn *websocket.Conn,
	hub *Hub,
) *Client {
	return &Client{
		send: make(chan []byte),
		conn: conn,
		hub:  hub,
	}
}

func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (c *Client) Write() {
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		c.hub.Broadcast(message)
	}
}
