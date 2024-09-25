package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

func ZapLoggerMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Passes to the next handler and if there is no error then err will be nil.
		err := c.Next()

		stop := time.Now()
		latency := stop.Sub(start)

		logger.Info("Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", latency),
		)

		return err
	}
}
