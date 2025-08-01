package http_v1

import (
	"cw/internal/controller/websocket"
	"cw/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	service *service.Service
	hub     *websocket.Hub
}

func NewHandler(logger *zap.Logger, service *service.Service, hub *websocket.Hub) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
		hub:     hub,
	}
}

func (h *Handler) InitRoutes() *fiber.App {
	app := fiber.New()

	app.Use(loggingMiddleware(h.logger))
	app.Use(recover.New())

	h.initChatRoutes(app)

	return app
}
