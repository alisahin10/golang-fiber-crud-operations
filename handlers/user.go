package handlers

import (
	"Fiber/models"
	"Fiber/services"
	"Fiber/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	Service *services.UserService
	Logger  *zap.Logger
}

func NewUserHandler(service *services.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		Service: service,
		Logger:  logger,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		h.Logger.Error("Failed to parse body", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	// Service layer.
	if err := h.Service.CreateUser(&user); err != nil {
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	responseUser := utils.ToResponseUser(&user)
	return c.Status(fiber.StatusCreated).JSON(responseUser)

}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.Service.GetUserByID(id)
	if err != nil {
		return utils.JSONErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	responseUser := utils.ToResponseUser(user)
	h.Logger.Info("User retrieved successfully", zap.String("user_id", user.ID))
	return c.JSON(responseUser)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	responseUsers := utils.ToResponseUsers(users)
	h.Logger.Info("All users retrieved successfully", zap.Int("count", len(users)))
	return c.JSON(responseUsers)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		h.Logger.Error("Failed to parse body", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := h.Service.UpdateUser(id, &updateData); err != nil {
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.Service.DeleteUser(id); err != nil {
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	// Retrieve both "name" and "email" query parameters
	nameQuery := c.Query("name")
	emailQuery := c.Query("email")

	// Log the search operation with the provided query parameters
	h.Logger.Info("Searching users", zap.String("name_query", nameQuery), zap.String("email_query", emailQuery))

	// If both queries are empty, return a bad request response
	if nameQuery == "" && emailQuery == "" {
		return utils.JSONErrorResponse(c, fiber.StatusBadRequest, "At least one query parameter (name or email) is required")
	}

	// Call the service to perform the search based on name and/or email
	users, err := h.Service.SearchUsersByNameOrEmail(nameQuery, emailQuery)
	if err != nil {
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, "Failed to search users")
	}
	return c.JSON(users)

}
