package http_v1

import (
	"cw/internal/controller/websocket"

	"github.com/gofiber/fiber/v2"
	fws "github.com/gofiber/websocket/v2"
)

func (h *Handler) initChatRoutes(app *fiber.App) {
	chat := app.Group("/chat")

	chat.Get("/", h.getChats)
	chat.Post("/", h.createChat)

	chat.Use("/ws", func(c *fiber.Ctx) error {
		if fws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	chat.Get("/ws", fws.New(h.serveWS))
}

func (h *Handler) serveWS(c *fws.Conn) {
	defer c.Close()

	userID := c.Query("id")
	client := websocket.NewClient(userID, c, h.hub)

	h.hub.Register(client)

	go client.Write()
	client.Read()
}

func (h *Handler) createChat(c *fiber.Ctx) error {
	var chatDTO struct {
		UserOneID string `json:"user_one_id"`
		UserTwoID string `json:"user_two_id"`
	}
	if err := c.BodyParser(&chatDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	return c.SendString("Chat endpoint")
}

func (h *Handler) getChats(c *fiber.Ctx) error {
	return c.SendString("Get chats endpoint")
}
