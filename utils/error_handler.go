package utils

import "github.com/gofiber/fiber/v2"

// JSONErrorResponse standardizes error responses in JSON format.
func JSONErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}
