package http_v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func loggingMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Since(start)

		logger.Info(
			"HTTP Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", stop),
		)

		return err
	}
}
