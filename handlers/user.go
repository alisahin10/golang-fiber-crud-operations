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

	// Delegate to the service layer.
	if err := h.Service.CreateUser(&user); err != nil {
		h.Logger.Error("Failed to create user", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	responseUser := models.ResponseUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	h.Logger.Info("User created successfully", zap.String("user_id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(responseUser)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.Service.GetUserByID(id)
	if err != nil {
		h.Logger.Error("Failed to get user by ID", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	responseUser := models.ResponseUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	h.Logger.Info("User retrieved successfully", zap.String("user_id", user.ID))
	return c.JSON(responseUser)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		h.Logger.Error("Failed to get all users", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	var responseUsers []models.ResponseUser
	for _, user := range users {
		responseUsers = append(responseUsers, models.ResponseUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

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
		h.Logger.Error("Failed to update user", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	h.Logger.Info("User updated successfully", zap.String("user_id", id))
	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.Service.DeleteUser(id); err != nil {
		h.Logger.Error("Failed to delete user", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	h.Logger.Info("User deleted successfully", zap.String("user_id", id))
	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

/*
func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("name")

	users, err := h.Service.SearchUsers(query)
	if err != nil {
		h.Logger.Error("Failed to search users", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if len(users) == 0 {
		h.Logger.Warn("No users found for query", zap.String("query", query))
		return utils.JSONErrorResponse(c, fiber.StatusNotFound, "No users found")
	}

	h.Logger.Info("Users found", zap.Int("count", len(users)))
	return c.JSON(users)
}

*/
/*
func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	// Get the search query from the query parameter
	query := c.Query("q")

	// Call the service to perform the search
	users, err := h.Service.SearchUsers(query)
	if err != nil {
		h.Logger.Error("Failed to search users", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// If no users are found, return a 404
	if len(users) == 0 {
		h.Logger.Warn("No users found for query", zap.String("query", query))
		return utils.JSONErrorResponse(c, fiber.StatusNotFound, "No users found")
	}

	h.Logger.Info("Users found", zap.Int("count", len(users)))
	return c.JSON(users)
}
*/

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
		h.Logger.Error("Failed to search users", zap.Error(err))
		return utils.JSONErrorResponse(c, fiber.StatusInternalServerError, "Failed to search users")
	}

	// If no users are found, return no content
	if len(users) == 0 {
		h.Logger.Info("No users found matching the criteria", zap.String("name_query", nameQuery), zap.String("email_query", emailQuery))
		return utils.JSONErrorResponse(c, fiber.StatusNoContent, "No users found matching the criteria")
	}

	h.Logger.Info("Users found", zap.Int("count", len(users)))
	return c.JSON(users)
}
