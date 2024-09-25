package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

func ZapLoggerMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details after processing
		logger.Info("Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", string(c.Request().Header.UserAgent())),
			zap.Duration("latency", time.Since(start)),
		)

		// If there's an error, log it
		if err != nil {
			logger.Error("Request error", zap.Error(err))
		}

		return err
	}
}
